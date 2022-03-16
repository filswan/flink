package service

import (
	"encoding/json"
	"flink-data/models"
	"fmt"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/logs"
	libutils "github.com/filswan/go-swan-lib/utils"
	"github.com/shopspring/decimal"
)

const JSON_RPC_VERSION = "2.0"
const JSON_RPC_ID = 1

type JsonRpcParams struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type FilscanDeal struct {
	Epoch                int64  `json:"epoch"`
	Label                string `json:"label"`
	Cid                  string `json:"cid"`
	DealId               int64  `json:"dealid"`
	Client               string `json:"client"`
	StartEpoch           int64  `json:"start_epoch"`
	EndEpoch             int64  `json:"end_epoch"`
	PieceCid             string `json:"piece_cid"`
	Provider             string `json:"provider"`
	PieceSize            string `json:"piece_size"`
	VerifiedDeal         bool   `json:"verified_deal"`
	ClientCollateral     string `json:"client_collateral"`
	ProviderCollateral   string `json:"provider_collateral"`
	StoragePricePerEpoch string `json:"storage_price_per_epoch"`
	BlockTime            int64  `json:"block_time"`
	ServiceStartTime     int64  `json:"service_start_time"`
	ServiceEndTime       int64  `json:"service_end_time"`
	Tag                  struct {
		TagCn   string `json:"tag_cn"`
		TagEn   string `json:"tag_en"`
		IsValid int    `json:"is_valid"`
	} `json:"tag"`
}

type JsonRpcResult struct {
	Id      int           `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Error   *JsonRpcError `json:"error"`
	Result  interface{}   `json:"result"`
}

type FilScanDealResult struct {
	JsonRpcResult
	Deal FilscanDeal `json:"result"`
}

type FilScanDealsResult struct {
	JsonRpcResult
	Result struct {
		Total int64          `json:"total"`
		Deals []*FilscanDeal `json:"deals"`
	} `json:"result"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

func GetFromJsonRpcApi(apiUrl, method string, params []interface{}, result interface{}) error {
	jsonRpcParams := JsonRpcParams{
		JsonRpc: JSON_RPC_VERSION,
		Method:  method,
		Params:  params,
		Id:      JSON_RPC_ID,
	}

	response, err := web.HttpGetNoToken(apiUrl, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	err = json.Unmarshal(response, result)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	logs.GetLogger().Info(result)
	return nil
}

func ConvertDeal2ChainLinkDeal(network models.Network, filscanDeal *FilscanDeal) *models.ChainLinkDeal {
	chainLinkDeal := models.ChainLinkDeal{
		NetworkId: network.Id,
	}

	chainLinkDeal.DealId = filscanDeal.DealId
	//chainLinkDeal.DealCid = filscanDeal.Cid
	chainLinkDeal.MessageCid = filscanDeal.Cid
	chainLinkDeal.Height = filscanDeal.Epoch
	chainLinkDeal.PieceCid = filscanDeal.PieceCid
	chainLinkDeal.VerifiedDeal = filscanDeal.VerifiedDeal

	storagePricePerEpoch, err := decimal.NewFromString(libutils.ConvertPrice2AttoFil(filscanDeal.StoragePricePerEpoch))
	if err != nil {
		logs.GetLogger().Error(err)
		chainLinkDeal.StoragePricePerEpoch = -1
	} else {
		chainLinkDeal.StoragePricePerEpoch = storagePricePerEpoch.BigInt().Int64()
	}

	//chainLinkDeal.Signature = deal.Signature
	//chainLinkDeal.SignatureType = deal.SignatureType
	chainLinkDeal.PieceSize = filscanDeal.PieceSize
	chainLinkDeal.StartHeight = filscanDeal.StartEpoch
	chainLinkDeal.EndHeight = filscanDeal.EndEpoch
	chainLinkDeal.Client = filscanDeal.Client
	chainLinkDeal.ClientCollateralFormat = libutils.GetPriceFormat("0 FIL")
	chainLinkDeal.Provider = filscanDeal.Provider
	//chainLinkDeal.ProviderTag = deal.ProviderTag
	//chainLinkDeal.VerifiedProvider = deal.VerifiedProvider
	chainLinkDeal.ProviderCollateralFormat = libutils.GetPriceFormat("0 FIL")
	chainLinkDeal.Status = filscanDeal.Tag.IsValid

	duration := chainLinkDeal.EndHeight - chainLinkDeal.StartHeight
	chainLinkDeal.StoragePrice = chainLinkDeal.StoragePricePerEpoch * duration

	chainLinkDeal.CreatedAt = filscanDeal.BlockTime
	logs.GetLogger().Info(chainLinkDeal)

	return &chainLinkDeal
}

func GetDealsOnDemand(dealId int64, networkName string) (*models.ChainLinkDeal, error) {
	network, err := models.GetNetworkByName(networkName)
	if err != nil {
		logs.GetLogger().Error()
		return nil, err
	}

	logs.GetLogger().Info("on demand requesting for:", dealId)

	chainLinkDeal, err := GetDealFromFilScanNetwork(*network, dealId)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	chainLinkDeals := []*models.ChainLinkDeal{}
	chainLinkDeals = append(chainLinkDeals, chainLinkDeal)

	err = models.AddChainLinkDeals(chainLinkDeals)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	logs.GetLogger().Info("inserted successfully on demand into db ,deal id:", dealId)

	return chainLinkDeal, nil
}

func GetMaxDealIdFromFilScanNetwork(network *models.Network) (*int64, error) {
	var params []interface{}
	params = append(params, "")
	params = append(params, 0)
	params = append(params, 1)

	filScanDealsResult := &FilScanDealsResult{}
	err := GetFromJsonRpcApi(network.ApiUrl, "filscan.GetMarketDeal", params, filScanDealsResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if filScanDealsResult.Error != nil {
		err := fmt.Errorf("error code:%d message:%s", filScanDealsResult.Error.Code, filScanDealsResult.Error.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &filScanDealsResult.Result.Deals[0].DealId, nil
}

func GetDealFromFilScanNetwork(network models.Network, dealId int64) (*models.ChainLinkDeal, error) {
	var params []interface{}
	params = append(params, dealId)

	filScanDealResult := &FilScanDealResult{}
	err := GetFromJsonRpcApi(network.ApiUrl, "filscan.GetMarketDealById", params, filScanDealResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if filScanDealResult.Error != nil {
		err := fmt.Errorf("error code:%d message:%s", filScanDealResult.Error.Code, filScanDealResult.Error.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	chainLinkDeal := ConvertDeal2ChainLinkDeal(network, &filScanDealResult.Deal)

	return chainLinkDeal, nil
}
