package service

import (
	"filink/data/models"
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

	deal := models.ChainLinkDealBase{
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
		PieceSizeFormat:          dealInternal.PieceSizeFormat,
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
	}

	return &deal, nil
}
