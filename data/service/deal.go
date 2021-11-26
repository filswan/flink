package service

import (
	"filink/data/models"
)

func GetDealById(dealId int64) (*models.ChainLinkDeal, error) {
	deal, err := models.GetDealById(dealId)
	if err != nil {
		//logs.GetLogger().Error(err)
		return nil, err
	}

	network, err := models.GetNetworkById(*deal.NetworkId)
	if err != nil {
		//logs.GetLogger().Error(err)
		return nil, err
	}
	deal.NetworkId = nil
	deal.NetworkName = network.Name

	return deal, nil
}
