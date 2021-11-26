package models

import (
	"filink/data/database"

	"github.com/filswan/go-swan-lib/logs"
)

type Network struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	ApiUrlPrefix string `json:"api_url_prefix"`
	Description  string `json:"description"`
}

func GetNetworkById(networkId int64) (*Network, error) {
	network := Network{}
	err := database.GetDB().Where("id=?", networkId).First(&network).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &network, nil
}

func GetNetworkByName(networkName string) (*Network, error) {
	network := Network{}
	err := database.GetDB().Where("name=?", networkName).First(&network).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &network, nil
}

func AddNetwork(network *Network) error {
	err := database.GetDB().Create(network).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}
