package service

import (
	"encoding/json"
	"fmt"
	"go-chainlink-data/common/constants"
	"go-chainlink-data/models"
	"strconv"
	"time"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func GetDealFromCalibration() error {
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
	//logs.GetLogger().Info(network.ApiUrlPrefix)
	for i := maxDealId + 1; ; i++ {
		apiUrlDeal := utils.UrlJoin(network.ApiUrlPrefix, strconv.FormatInt(i, 10))
		response := client.HttpGetNoToken(apiUrlDeal, nil)
		if response == "" {
			err := fmt.Errorf("no response from:%s", apiUrlDeal)
			//logs.GetLogger().Error(err)
			return err
		}

		chainLinkDealCalibrationResult := &models.ChainLinkDealCalibrationResult{}
		err := json.Unmarshal([]byte(response), chainLinkDealCalibrationResult)
		if err != nil {
			//logs.GetLogger().Error(err)
			return err
		}

		//logs.GetLogger().Info(apiUrlDeal, ",", chainLinkDeal.Code, ",", chainLinkDeal.Message, ",", chainLinkDeal.Message)

		currentMilliSec := time.Now().UnixNano() / 1e6

		if currentMilliSec-lastInsertAt >= constants.BULK_INSERT_INTERVAL_MILLI_SEC && len(chainLinkDeals) >= 1 {
			err := models.AddChainLinkDeals(chainLinkDeals)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			chainLinkDeals = []*models.ChainLinkDeal{}
		}

		if chainLinkDealCalibrationResult.Code == constants.CALIBRATION_DEAL_NOT_FOUND {
			continue
		}

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

		chainLinkDeals = append(chainLinkDeals, &chainLinkDeal)
		if len(chainLinkDeals) >= constants.BULK_INSERT_CHAINLINK_LIMIT {
			err := models.AddChainLinkDeals(chainLinkDeals)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			chainLinkDeals = []*models.ChainLinkDeal{}
		}
		logs.GetLogger().Info(chainLinkDeal)
	}
}
