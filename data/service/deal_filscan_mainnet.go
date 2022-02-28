package service

import (
	"encoding/json"
	"fmt"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/logs"
)

const JSON_RPC_VERSION = "2.0"
const JSON_RPC_ID = 1
const JSON_RPC_METHOD_FILSCAN_GET_MARKET_DEAL_BY_ID = ""

const API_URL = "https://api.filscan.io:8700/rpc/v1"

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
}

type JsonRpcResult struct {
	Id      int           `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Error   *JsonRpcError `json:"error"`
}

type FilScanDealResult struct {
	JsonRpcResult
	Deal FilscanDeal `json:"result"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GetDealFromFilScanMainNet(dealId int64) (*FilscanDeal, error) {
	var params []interface{}
	params = append(params, dealId)

	jsonRpcParams := JsonRpcParams{
		JsonRpc: JSON_RPC_VERSION,
		Method:  JSON_RPC_METHOD_FILSCAN_GET_MARKET_DEAL_BY_ID,
		Params:  params,
		Id:      JSON_RPC_ID,
	}

	response, err := web.HttpGetNoToken(API_URL, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	filScanDealResult := &FilScanDealResult{}
	err = json.Unmarshal(response, filScanDealResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if filScanDealResult.Error != nil {
		err := fmt.Errorf("error code:%d message:%s", filScanDealResult.Error.Code, filScanDealResult.Error.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &filScanDealResult.Deal, nil
}
