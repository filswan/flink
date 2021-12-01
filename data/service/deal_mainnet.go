package service

import (
	"encoding/json"
	"flink-data/common/constants"
	"flink-data/config"
	"flink-data/models"
	"fmt"
	"strconv"
	"time"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	libutils "github.com/filswan/go-swan-lib/utils"
	"github.com/shopspring/decimal"
)

func GetDealsFromMainnetLoop() {
	for {
		logs.GetLogger().Info("start")
		err := GetDealsFromMainnet()
		if err != nil {
			logs.GetLogger().Error()
		}

		logs.GetLogger().Info("sleep")
		time.Sleep(time.Minute * 1)
	}
}

func GetDealsFromMainnet() error {
	network, err := models.GetNetworkByName(constants.NETWORK_MAINNET)
	if err != nil {
		logs.GetLogger().Error()
		return err
	}

	maxDealId, err := models.GetMaxDealId(network.Id)
	if err != nil {
		logs.GetLogger().Error()
		return err
	}

	chainLinkDeals := []*models.ChainLinkDeal{}

	logs.GetLogger().Info("max deal id last scanned:", maxDealId)

	lastInsertAt := time.Now().UnixNano() / 1e6

	startDealId := maxDealId + 1
	lastDealId := maxDealId

	//logs.GetLogger().Info(network.ApiUrlPrefix)

	bulkInsertChainLinkLimit := config.GetConfig().ChainLink.BulkInsertChainlinkLimit
	bulkInsertIntervalMilliSec := config.GetConfig().ChainLink.BulkInsertIntervalMilliSec
	dealIdMaxInterval := config.GetConfig().ChainLink.DealIdIntervalMax

	logs.GetLogger().Info("scanned from:", startDealId)
	for i := startDealId; ; i++ {
		foundDeal := false
		chainLinkDeal, err := GetDealFromMainnet(*network, i)
		if err != nil {
			logs.GetLogger().Error(err)
		} else {
			chainLinkDeals = append(chainLinkDeals, chainLinkDeal)
			lastDealId = chainLinkDeal.DealId
			foundDeal = true
		}

		dealIdInterval := i - lastDealId
		//logs.GetLogger().Info(dealIdInterval)
		currentMilliSec := time.Now().UnixNano() / 1e6
		if len(chainLinkDeals) >= bulkInsertChainLinkLimit ||
			(currentMilliSec-lastInsertAt >= bulkInsertIntervalMilliSec && len(chainLinkDeals) >= 1) ||
			(dealIdInterval > dealIdMaxInterval && len(chainLinkDeals) >= 1) ||
			(!foundDeal && len(chainLinkDeals) >= 1) {
			logs.GetLogger().Info("insert into db, deals count:", len(chainLinkDeals), ",deal id interval:", dealIdInterval, ",last insert at:", lastInsertAt, ",current milliseconds:", currentMilliSec)
			err := models.AddChainLinkDeals(chainLinkDeals)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			chainLinkDeals = []*models.ChainLinkDeal{}
			lastInsertAt = currentMilliSec
		}

		if dealIdInterval >= dealIdMaxInterval {
			logs.GetLogger().Info("last deal id scanned:", i, ",scanned from:", startDealId)
			return nil
		}
	}
}

type MainnetHeightResult struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    CalibrationHeight `json:"data"`
}

type MainnetHeight struct {
	TipSetHeight int64 `json:"tipSetHeight"`
	CountDown    int   `json:"countDown"`
}

func GetHeightFromMainnet(network models.Network) (int64, error) {
	response := client.HttpGetNoToken(network.ApiUrlHeight, nil)
	if response == "" {
		err := fmt.Errorf("no response from:%s", network.ApiUrlHeight)
		//logs.GetLogger().Error(err)
		return -1, err
	}

	mainnetHeightResult := &MainnetHeightResult{}
	err := json.Unmarshal([]byte(response), mainnetHeightResult)
	if err != nil {
		err := fmt.Errorf("%s from:%s", err.Error(), network.ApiUrlHeight)
		//logs.GetLogger().Error(err)
		return -1, err
	}

	if mainnetHeightResult.Code != 200 {
		err := fmt.Errorf("code:%d,message:%s", mainnetHeightResult.Code, mainnetHeightResult.Message)
		return -1, err
	}

	return mainnetHeightResult.Data.TipSetHeight, nil
}

func GetDealFromMainnet(network models.Network, dealId int64) (*models.ChainLinkDeal, error) {
	apiUrlDeal := libutils.UrlJoin(network.ApiUrlPrefix, strconv.FormatInt(dealId, 10))
	response := client.HttpGetNoToken(apiUrlDeal, nil)
	if response == "" {
		err := fmt.Errorf("deal_id:%d,no response from:%s", dealId, apiUrlDeal)
		//logs.GetLogger().Error(err)
		return nil, err
	}

	chainLinkDealMainnetResult := &models.ChainLinkDealMainnetResult{}
	err := json.Unmarshal([]byte(response), chainLinkDealMainnetResult)
	if err != nil {
		err := fmt.Errorf("deal_id:%d,%s", dealId, err.Error())
		//logs.GetLogger().Error(err)
		return nil, err
	}

	if chainLinkDealMainnetResult.Code == constants.CALIBRATION_DEAL_NOT_FOUND {
		err := fmt.Errorf("deal_id:%d,code:%d,message:%s", dealId, chainLinkDealMainnetResult.Code, chainLinkDealMainnetResult.Message)
		return nil, err
	}
	//logs.GetLogger().Info(apiUrlDeal, ",", chainLinkDeal.Code, ",", chainLinkDeal.Message, ",", chainLinkDeal.Message)

	deal := chainLinkDealMainnetResult.Data
	chainLinkDeal := models.ChainLinkDeal{
		NetworkId: network.Id,
	}

	chainLinkDeal.DealId = deal.DealId
	chainLinkDeal.DealCid = deal.DealCid
	chainLinkDeal.MessageCid = deal.MessageCid
	chainLinkDeal.Height = deal.Height
	chainLinkDeal.PieceCid = deal.PieceCid
	chainLinkDeal.VerifiedDeal = deal.VerifiedDeal

	storagePricePerEpoch, err := decimal.NewFromString(libutils.ConvertPrice2AttoFil(deal.StoragePricePerEpoch))
	if err != nil {
		logs.GetLogger().Error(err)
		chainLinkDeal.StoragePricePerEpoch = -1
	} else {
		chainLinkDeal.StoragePricePerEpoch = storagePricePerEpoch.BigInt().Int64()
	}

	chainLinkDeal.Signature = deal.Signature
	chainLinkDeal.SignatureType = deal.SignatureType
	chainLinkDeal.PieceSizeFormat = deal.PieceSizeFormat
	chainLinkDeal.StartHeight = deal.StartHeight
	chainLinkDeal.EndHeight = deal.EndHeight
	chainLinkDeal.Client = deal.Client
	chainLinkDeal.ClientCollateralFormat = libutils.GetPriceFormat("0 FIL")
	chainLinkDeal.Provider = deal.Provider
	chainLinkDeal.ProviderTag = deal.ProviderTag
	chainLinkDeal.VerifiedProvider = deal.VerifiedProvider
	chainLinkDeal.ProviderCollateralFormat = libutils.GetPriceFormat("0 FIL")
	chainLinkDeal.Status = deal.Status

	duration := chainLinkDeal.EndHeight - chainLinkDeal.StartHeight
	chainLinkDeal.StoragePrice = chainLinkDeal.StoragePricePerEpoch * duration

	timeT, err := time.Parse("2006-01-02 15:04:05", deal.CreatedAt)
	if err != nil {
		logs.GetLogger().Error(err)
	} else {
		chainLinkDeal.CreatedAt = timeT.UnixNano() / 1e9
	}
	logs.GetLogger().Info(chainLinkDeal)

	return &chainLinkDeal, nil
}