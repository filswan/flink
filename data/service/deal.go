package service

import (
	"encoding/json"
	"flink-data/common/constants"
	"flink-data/models"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/logs"
	libutils "github.com/filswan/go-swan-lib/utils"
)

func GetDealById(dealId int64, networkName string) (*models.ChainLinkDealBase, error) {
	network, err := models.GetNetworkByName(networkName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	dealInternal, err := models.GetDealByIdAndNetwork(dealId, int(network.Id))
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if dealInternal == nil {
		if network.Name == constants.NETWORK_CALIBRATION {
			dealInternal, err = GetDealOnDemandFromCalibration(dealId)
			if err != nil {
				logs.GetLogger().Error(err)
				return nil, err
			}
		}

		if network.Name == constants.NETWORK_MAINNET {
			dealInternal, err = GetDealsOnDemandFromMainnet(dealId)
			if err != nil {
				logs.GetLogger().Error(err)
				return nil, err
			}
		}
	}

	if dealInternal != nil {
		dealInternal.NetworkName = network.Name

		deal := models.GetChainLinkDealBase(*dealInternal)

		return &deal, nil
	} else {
		return nil, nil
	}

}

func GetLatestDealByNetwork(networkId int64) (*models.ChainLinkDealBase, error) {
	maxDealId, err := models.GetMaxDealId(networkId)

	if err != nil {
		return nil, err
	}

	dealInternal, err := models.GetDealByIdAndNetwork(maxDealId, int(networkId))
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	network, err := models.GetNetworkById(dealInternal.NetworkId)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	dealInternal.NetworkName = network.Name

	deal := models.GetChainLinkDealBase(*dealInternal)

	return &deal, nil
}

func GetCurrentMaxDealFromChainLink(network models.Network) (int64, error) {
	apiUrlDeal := libutils.UrlJoin(network.ApiUrlStorage)
	response, err := web.HttpPostNoToken(apiUrlDeal, map[string]interface{}{
		"client":    "",
		"keyword":   "",
		"pageIndex": 1,
		"pageSize":  1,
		"provider":  "",
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return -1, err
	}

	if network.Name == constants.NETWORK_CALIBRATION {
		chainLinkDealCalibrationResult := &models.ChainLinkDealCalibrationArrayResult{}
		err = json.Unmarshal(response, chainLinkDealCalibrationResult)
		if err != nil {
			logs.GetLogger().Error(err)
			return -1, err
		}

		if chainLinkDealCalibrationResult.Code == constants.CALIBRATION_DEAL_FOUND {
			return chainLinkDealCalibrationResult.Data[0].DealId, nil
		}
	} else {
		chainLinkDealMainnetResult := &models.ChainLinkDealMainnetArrayResult{}
		err = json.Unmarshal(response, chainLinkDealMainnetResult)
		if err != nil {
			logs.GetLogger().Error(err)
			return -1, err
		}

		if chainLinkDealMainnetResult.Code == constants.MAINNET_DEAL_FOUND {
			return chainLinkDealMainnetResult.Data[0].DealId, nil
		}
	}

	return -1, nil
}
