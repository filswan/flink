package service

import (
	"flink-data/common/constants"
	"flink-data/config"
	"flink-data/models"
	"time"

	"github.com/filswan/go-swan-lib/logs"
)

func GetDealsFromCalibrationLoop() {
	for {
		logs.GetLogger().Info("start")

		network, err := models.GetNetworkByName(constants.NETWORK_CALIBRATION)
		if err != nil {
			logs.GetLogger().Error()
			return
		}

		maxDealIdOnFilScan, err := GetMaxDealIdFromFilScanNetwork(network)
		if err != nil {
			logs.GetLogger().Error()
			return
		}

		err = GetDealsFromCalibration(network, *maxDealIdOnFilScan)
		if err != nil {
			logs.GetLogger().Error()
		}

		logs.GetLogger().Info("sleep")
		time.Sleep(time.Minute * 1)
	}
}

func GetDealsFromCalibration(network *models.Network, maxDealIdOnFilScan int64) error {
	maxDealId, err := models.GetMaxDealId(network.Id)
	if err != nil {
		logs.GetLogger().Error()
		return err
	}

	chainLinkDeals := []*models.ChainLinkDeal{}

	logs.GetLogger().Info("max deal id last scanned:", maxDealId, " scanned from:", maxDealId+1)

	bulkInsertChainLinkLimit := config.GetConfig().ChainLink.BulkInsertChainlinkLimit
	bulkInsertIntervalMilliSec := config.GetConfig().ChainLink.BulkInsertIntervalMilliSec
	lastInsertAt := time.Now().UnixNano() / 1e6

	for i := maxDealId + 1; i <= maxDealIdOnFilScan; i++ {
		chainLinkDeal, err := GetDealFromFilScanNetwork(*network, i)
		if err != nil {
			logs.GetLogger().Error("deal id:", i, " ", err)
		} else {
			chainLinkDeals = append(chainLinkDeals, chainLinkDeal)
		}

		//logs.GetLogger().Info(dealIdInterval)
		currentMilliSec := time.Now().UnixNano() / 1e6
		if len(chainLinkDeals) >= bulkInsertChainLinkLimit || i == maxDealIdOnFilScan ||
			(currentMilliSec-lastInsertAt >= bulkInsertIntervalMilliSec && len(chainLinkDeals) >= 1) {
			logs.GetLogger().Info("insert into db, deals count:", len(chainLinkDeals), ",last insert at:", lastInsertAt, ",current milliseconds:", currentMilliSec)
			err := models.AddChainLinkDeals(chainLinkDeals)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			chainLinkDeals = []*models.ChainLinkDeal{}
			lastInsertAt = currentMilliSec
		}
	}
	return nil
}
