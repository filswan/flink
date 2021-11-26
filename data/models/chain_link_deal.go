package models

import (
	"filink/data/database"
	"fmt"
	"strings"

	"github.com/filswan/go-swan-lib/logs"
)

type ChainLinkDeal struct {
	DealId                   int64  `json:"deal_id"`
	DealCid                  string `json:"deal_cid"`
	MessageCid               string `json:"message_cid"`
	Height                   int64  `json:"height"`
	PieceCid                 string `json:"piece_cid"`
	VerifiedDeal             bool   `json:"verified_deal"`
	StoragePricePerEpoch     string `json:"storage_price_per_epoch"`
	Signature                string `json:"signature"`
	SignatureType            string `json:"signature_type"`
	CreatedAtSrc             string `json:"created_at_src"`
	CreatedAt                int64  `json:"created_at"`
	PieceSizeFormat          string `json:"piece_size_format"`
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
}

type ChainLinkDealInternal struct {
	ChainLinkDeal
	NetworkId int64 `json:"network_id"`
}

func AddChainLinkDeal(chainLinkDeal *ChainLinkDealInternal) error {
	err := database.GetDB().Create(chainLinkDeal).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

func AddChainLinkDeals(chainLinkDeals []*ChainLinkDealInternal) error {
	if len(chainLinkDeals) <= 0 {
		err := fmt.Errorf("no deal in chainLinkDeals")
		return err
	}

	sql := "insert into chain_link_deal (deal_id,network_id,deal_cid,message_cid,height,piece_cid,verified_deal,storage_price_per_epoch,signature,signature_type,created_at_src,created_at,"
	sql = sql + "piece_size_format,start_height,end_height,client,client_collateral_format,provider,provider_tag,verified_provider,provider_collateral_format,status) values"
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
		valueArgs = append(valueArgs, deal.CreatedAtSrc)
		valueArgs = append(valueArgs, deal.CreatedAt)
		valueArgs = append(valueArgs, deal.PieceSizeFormat)
		valueArgs = append(valueArgs, deal.StartHeight)
		valueArgs = append(valueArgs, deal.EndHeight)
		valueArgs = append(valueArgs, deal.Client)
		valueArgs = append(valueArgs, deal.ClientCollateralFormat)
		valueArgs = append(valueArgs, deal.Provider)
		valueArgs = append(valueArgs, deal.ProviderTag)
		valueArgs = append(valueArgs, deal.VerifiedProvider)
		valueArgs = append(valueArgs, deal.ProviderCollateralFormat)
		valueArgs = append(valueArgs, deal.Status)
	}

	sql = fmt.Sprintf("%s %s", sql, strings.Join(valueStrings, ","))

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
	var chainLinkDeal ChainLinkDealInternal
	err := database.GetDB().Raw(sql, networkId).Scan(&chainLinkDeal).Error
	if err != nil {
		logs.GetLogger().Error(err)
		return 0, err
	}
	return chainLinkDeal.DealId, nil
}

func GetDealById(dealId int64) (*ChainLinkDealInternal, error) {
	chainLinkDeal := ChainLinkDealInternal{}
	err := database.GetDB().Where("deal_id=?", dealId).First(&chainLinkDeal).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &chainLinkDeal, nil
}
