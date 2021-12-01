package models

type ChainLinkDealMainnet struct {
	DealId                   int64  `json:"DealId"`
	DealCid                  string `json:"dealCid"`
	MessageCid               string `json:"messageCid"`
	Height                   int64  `json:"height"`
	PieceCid                 string `json:"pieceCid"`
	VerifiedDeal             bool   `json:"verifiedDeal"`
	StoragePricePerEpoch     string `json:"storagePricePerEpoch"`
	Signature                string `json:"signature"`
	SignatureType            string `json:"signatureType"`
	CreatedAt                string `json:"createdAt"`
	PieceSizeFormat          string `json:"pieceSizeFormat"`
	StartHeight              int64  `json:"satrtHeight"`
	EndHeight                int64  `json:"endHeight"`
	Client                   string `json:"client"`
	ClientCollateralFormat   string `json:"clientCollateralFormat"`
	Provider                 string `json:"provider"`
	ProviderTag              string `json:"providerTag"`
	VerifiedProvider         int    `json:"providerIsVerified"`
	ProviderCollateralFormat string `json:"providerCollateralFormat"`
	Status                   int    `json:"status"`
}

type ChainLinkDealMainnetResult struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    ChainLinkDealMainnet `json:"data"`
}
