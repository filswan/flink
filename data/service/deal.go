package service

import (
	"flink-data/models"

	"github.com/filswan/go-swan-lib/logs"
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
		dealInternal, err = GetDealsOnDemand(dealId, network.Name)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
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
	dealInternal, err := models.GetLastDeal(networkId)
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
