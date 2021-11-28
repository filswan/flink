package service

import (
	"flink-data/models"
)

func GetDealById(dealId int64) (*models.ChainLinkDealBase, error) {
	dealInternal, err := models.GetDealById(dealId)
	if err != nil {
		//logs.GetLogger().Error(err)
		return nil, err
	}

	network, err := models.GetNetworkById(dealInternal.NetworkId)
	if err != nil {
		//logs.GetLogger().Error(err)
		return nil, err
	}
	dealInternal.NetworkName = network.Name

	deal := models.GetChainLinkDealBase(*dealInternal)

	return &deal, nil
}
