package models

import (
	"flink-data/database"
	"fmt"
	"strings"

	"github.com/filswan/go-swan-lib/logs"
)

type ChainLinkDealBase struct {
	DealId                   int64  `json:"deal_id"`
	DealCid                  string `json:"deal_cid"`
	MessageCid               string `json:"message_cid"`
	Height                   int64  `json:"height"`
	PieceCid                 string `json:"piece_cid"`
	VerifiedDeal             bool   `json:"verified_deal"`
	StoragePricePerEpoch     int64  `json:"storage_price_per_epoch"`
	Signature                string `json:"signature"`
	SignatureType            string `json:"signature_type"`
	CreatedAt                int64  `json:"created_at"`
	PieceSize                string `json:"piece_size"`
	StartHeight              int64  `json:"start_height"`
	EndHeight                int64  `json:"end_height"`
	Client                   string `json:"client"`
	ClientCollateralFormat   string `json:"client_collateral_format"`
	Provider                 string `json:"provider"`
	ProviderTag              string `json:"provider_tag"`
	VerifiedProvider         int    `json:"verified_provider"`
	ProviderCollateralFormat string `json:"provider_collateral_format"`
	Status                   int    `json:"status"`
	NetworkName              string `json:"network_name"`
	StoragePrice             int64  `json:"storage_price"`
}

type ChainLinkDeal struct {
	ChainLinkDealBase
	NetworkId int64 `json:"network_id"`
}

func GetChainLinkDealBase(dealInternal ChainLinkDeal) ChainLinkDealBase {
	deal := ChainLinkDealBase{
		DealId:                   dealInternal.DealId,
		DealCid:                  dealInternal.DealCid,
		MessageCid:               dealInternal.MessageCid,
		Height:                   dealInternal.Height,
		PieceCid:                 dealInternal.PieceCid,
		VerifiedDeal:             dealInternal.VerifiedDeal,
		StoragePricePerEpoch:     dealInternal.StoragePricePerEpoch,
		Signature:                dealInternal.Signature,
		SignatureType:            dealInternal.SignatureType,
		CreatedAt:                dealInternal.CreatedAt,
		PieceSize:                dealInternal.PieceSize,
		StartHeight:              dealInternal.StartHeight,
		EndHeight:                dealInternal.EndHeight,
		Client:                   dealInternal.Client,
		ClientCollateralFormat:   dealInternal.ClientCollateralFormat,
		Provider:                 dealInternal.Provider,
		ProviderTag:              dealInternal.ProviderTag,
		VerifiedProvider:         dealInternal.VerifiedProvider,
		ProviderCollateralFormat: dealInternal.ProviderCollateralFormat,
		Status:                   dealInternal.Status,
		NetworkName:              dealInternal.NetworkName,
		StoragePrice:             dealInternal.StoragePrice,
	}

	return deal
}

func AddChainLinkDeal(chainLinkDeal *ChainLinkDeal) error {
	err := database.GetDB().Create(chainLinkDeal).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

func AddChainLinkDeals(chainLinkDeals []*ChainLinkDeal) error {
	if len(chainLinkDeals) <= 0 {
		err := fmt.Errorf("no deal in chainLinkDeals")
		return err
	}

	sql := "insert into chain_link_deal (deal_id,network_id,deal_cid,message_cid,height,piece_cid,verified_deal,storage_price_per_epoch,signature,signature_type,created_at,"
	sql = sql + "piece_size,start_height,end_height,client,client_collateral_format,provider,provider_tag,verified_provider,provider_collateral_format,status,storage_price) values"
	valueStrings := []string{}

	valueArgs := []interface{}{}
	for _, deal := range chainLinkDeals {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, deal.DealId)
		valueArgs = append(valueArgs, deal.NetworkId)
		valueArgs = append(valueArgs, deal.DealCid)
		valueArgs = append(valueArgs, deal.MessageCid)
		valueArgs = append(valueArgs, deal.Height)
		valueArgs = append(valueArgs, deal.PieceCid)
		valueArgs = append(valueArgs, deal.VerifiedDeal)
		valueArgs = append(valueArgs, deal.StoragePricePerEpoch)
		valueArgs = append(valueArgs, deal.Signature)
		valueArgs = append(valueArgs, deal.SignatureType)
		valueArgs = append(valueArgs, deal.CreatedAt)
		valueArgs = append(valueArgs, deal.PieceSize)
		valueArgs = append(valueArgs, deal.StartHeight)
		valueArgs = append(valueArgs, deal.EndHeight)
		valueArgs = append(valueArgs, deal.Client)
		valueArgs = append(valueArgs, deal.ClientCollateralFormat)
		valueArgs = append(valueArgs, deal.Provider)
		valueArgs = append(valueArgs, deal.ProviderTag)
		valueArgs = append(valueArgs, deal.VerifiedProvider)
		valueArgs = append(valueArgs, deal.ProviderCollateralFormat)
		valueArgs = append(valueArgs, deal.Status)
		valueArgs = append(valueArgs, deal.StoragePrice)
	}

	sql = fmt.Sprintf("%s %s", sql, strings.Join(valueStrings, ","))

	onDuplicateKey := "network_id=values(network_id),deal_cid=values(deal_cid),message_cid=values(message_cid),height=values(height),piece_cid=values(piece_cid),"
	onDuplicateKey = onDuplicateKey + "verified_deal=values(verified_deal),storage_price_per_epoch=values(storage_price_per_epoch),signature=values(signature),signature_type=values(signature_type),"
	onDuplicateKey = onDuplicateKey + "created_at=values(created_at),piece_size=values(piece_size),start_height=values(start_height),end_height=values(end_height),"
	onDuplicateKey = onDuplicateKey + "client=values(client),client_collateral_format=values(client_collateral_format),provider=values(provider),provider_tag=values(provider_tag),"
	onDuplicateKey = onDuplicateKey + "verified_provider=values(verified_provider),provider_collateral_format=values(provider_collateral_format),status=values(status),storage_price=values(storage_price)"
	sql = sql + " on duplicate key update " + onDuplicateKey

	err := database.GetDB().Exec(sql, valueArgs...).Error

	if err != nil {
		logs.GetLogger().Info(sql)
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

func GetMaxDealId(networkId int64) (int64, error) {
	sql := "select max(deal_id) deal_id from chain_link_deal where network_id=?"
	var chainLinkDeal ChainLinkDeal
	err := database.GetDB().Raw(sql, networkId).Scan(&chainLinkDeal).Error
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}
	return chainLinkDeal.DealId, nil
}

func GetDealById(dealId int64) (*ChainLinkDeal, error) {
	chainLinkDeal := ChainLinkDeal{}
	err := database.GetDB().Where("deal_id=?", dealId).First(&chainLinkDeal).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &chainLinkDeal, nil
}

func GetDealByIdAndNetwork(dealId int64, networkId int) (*ChainLinkDeal, error) {
	var chainLinkDeal []*ChainLinkDeal

	err := database.GetDB().Where("deal_id=? and network_id=?", dealId, networkId).Find(&chainLinkDeal).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if len(chainLinkDeal) > 0 {
		return chainLinkDeal[0], nil
	}

	return nil, nil
}

func GetLastDeal(networkId int64) (*ChainLinkDeal, error) {
	maxDealId, err := GetMaxDealId(networkId)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	lastDeal, err := GetDealById(maxDealId)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return lastDeal, nil
}
