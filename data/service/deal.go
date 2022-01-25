package service

import (
	"flink-data/models"

	"github.com/filswan/go-swan-lib/logs"
)

func GetDealById(dealId int64, networkId int) (*models.ChainLinkDealBase, error) {
	dealInternal, err := models.GetDealByIdAndNetwork(dealId, networkId)
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
