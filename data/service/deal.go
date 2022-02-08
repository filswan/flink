package service

import (
	"flink-data/common/constants"
	"flink-data/models"
	"fmt"

	"github.com/filswan/go-swan-lib/logs"
)

func GetDealById(dealId int64, networkId int) (*models.ChainLinkDealBase, error) {
	dealInternal, err := models.GetDealByIdAndNetwork(dealId, networkId)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	network, err := models.GetNetworkById(int64(networkId))

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
		err := fmt.Errorf("deal not found for %d in %s", dealId, network.Name)
		logs.GetLogger().Error(err)
		return nil, err
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
