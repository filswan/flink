package service

import (
	"encoding/json"
	"filecoin-data-provider/data/common/constants"
	"filecoin-data-provider/data/models"
	"fmt"
	"strconv"
	"time"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func GetDealsFromCalibrationLoop() {
	for {
		err := GetDealsFromCalibration()
		if err != nil {
			logs.GetLogger().Error()
		}

		time.Sleep(time.Minute * 1)
	}
}

func GetDealsFromCalibration() error {
	network, err := models.GetNetworkByName(constants.NETWORK_CALIBRATION)
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
	for i := startDealId; ; i++ {
		chainLinkDeal, err := GetDealFromCalibration(*network, i)
		if err != nil {
			logs.GetLogger().Error(err)
			continue
		} else {
			chainLinkDeals = append(chainLinkDeals, chainLinkDeal)
			lastDealId = chainLinkDeal.DealId
		}

		currentMilliSec := time.Now().UnixNano() / 1e6
		if len(chainLinkDeals) >= constants.BULK_INSERT_CHAINLINK_LIMIT || (currentMilliSec-lastInsertAt >= constants.BULK_INSERT_INTERVAL_MILLI_SEC && len(chainLinkDeals) >= 1) {
			err := models.AddChainLinkDeals(chainLinkDeals)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			chainLinkDeals = []*models.ChainLinkDeal{}
			lastInsertAt = currentMilliSec
		}

		if i-lastDealId > 10000 {
			err := fmt.Errorf("no deal for the last 10000 deal id")
			logs.GetLogger().Error(err)
			return err
		}
	}
}

func GetDealFromCalibration(network models.Network, dealId int64) (*models.ChainLinkDeal, error) {
	apiUrlDeal := utils.UrlJoin(network.ApiUrlPrefix, strconv.FormatInt(dealId, 10))
	response := client.HttpGetNoToken(apiUrlDeal, nil)
	if response == "" {
		err := fmt.Errorf("deal_id:%d,no response from:%s", dealId, apiUrlDeal)
		//logs.GetLogger().Error(err)
		return nil, err
	}

	chainLinkDealCalibrationResult := &models.ChainLinkDealCalibrationResult{}
	err := json.Unmarshal([]byte(response), chainLinkDealCalibrationResult)
	if err != nil {
		err := fmt.Errorf("deal_id:%d,%s", dealId, err.Error())
		//logs.GetLogger().Error(err)
		return nil, err
	}

	if chainLinkDealCalibrationResult.Code == constants.CALIBRATION_DEAL_NOT_FOUND {
		err := fmt.Errorf("deal_id:%d,code:%d,message:%s", dealId, chainLinkDealCalibrationResult.Code, chainLinkDealCalibrationResult.Message)
		return nil, err
	}
	//logs.GetLogger().Info(apiUrlDeal, ",", chainLinkDeal.Code, ",", chainLinkDeal.Message, ",", chainLinkDeal.Message)

	deal := chainLinkDealCalibrationResult.Data
	chainLinkDeal := models.ChainLinkDeal{
		DealId:                   deal.DealId,
		NetworkId:                network.Id,
		DealCid:                  deal.DealCid,
		MessageCid:               deal.MessageCid,
		Height:                   deal.Height,
		PieceCid:                 deal.PieceCid,
		VerifiedDeal:             deal.VerifiedDeal,
		StoragePricePerEpoch:     deal.StoragePricePerEpoch,
		Signature:                deal.Signature,
		SignatureType:            deal.SignatureType,
		CreatedAtSrc:             deal.CreatedAtSrc,
		PieceSizeFormat:          deal.PieceSizeFormat,
		StartHeight:              deal.StartHeight,
		EndHeight:                deal.EndHeight,
		Client:                   deal.Client,
		ClientCollateralFormat:   deal.ClientCollateralFormat,
		Provider:                 deal.Provider,
		ProviderTag:              deal.ProviderTag,
		VerifiedProvider:         deal.VerifiedProvider,
		ProviderCollateralFormat: deal.ProviderCollateralFormat,
		Status:                   deal.Status,
	}
	timeT, err := time.Parse("2006-01-02 15:04:05", chainLinkDeal.CreatedAtSrc)
	if err != nil {
		logs.GetLogger().Error(err)
	} else {
		chainLinkDeal.CreatedAt = timeT.UnixNano() / 1e6
	}
	logs.GetLogger().Info(chainLinkDeal)

	return &chainLinkDeal, nil
}
