package routers

import (
	"flink-data/common"
	"flink-data/models"
	"flink-data/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Deal(router *gin.RouterGroup) {
	router.POST("", GetDeal)
	router.POST("latest", GetLatestDeal)
}

func GetDeal(c *gin.Context) {
	var dealNetworkRequest_v1 DealNetworkRequest_v1
	var dealNetworkRequest = new(DealNetworkRequest)
	err := c.ShouldBindBodyWith(&dealNetworkRequest_v1, binding.JSON)
	if err != nil {
		err := c.ShouldBindBodyWith(&dealNetworkRequest, binding.JSON)
		if err != nil {
			logs.GetLogger().Error(err)
			c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
			return
		} else {
			dealNetworkRequest.DealId = dealNetworkRequest.DealId
			dealNetworkRequest.NetworkName = dealNetworkRequest_v1.NetworkName
		}
	} else {
		dealNetworkRequest.DealId = strconv.Itoa(dealNetworkRequest_v1.DealId)
		dealNetworkRequest.NetworkName = dealNetworkRequest_v1.NetworkName
	}

	dealIdStr := dealNetworkRequest.DealId
	dealId, err := strconv.ParseInt(dealIdStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("deal id must be numeric")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	networkName := dealNetworkRequest.NetworkName
	if networkName == "" {
		err := fmt.Errorf("network id must be provided")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	deal, err := service.GetDealById(dealId, networkName)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	if deal == nil {
		err := fmt.Errorf("deal not found")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	mapObject := map[string]interface{}{
		"deal": *deal,
	}

	c.JSON(http.StatusOK, common.CreateSuccessResponse(mapObject))
}

func AddMiner(c *gin.Context) {
}

func GetLatestDeal(c *gin.Context) {
	var dealNetworkRequest DealNetworkRequest
	err := c.BindJSON(&dealNetworkRequest)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	networkName := dealNetworkRequest.NetworkName
	if networkName == "" {
		err := fmt.Errorf("network id must be provided")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	network, err := models.GetNetworkByName(networkName)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	if network == nil {
		err := fmt.Errorf("network is not valid")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	deal, err := service.GetLatestDealByNetwork(network.Id)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	mapObject := map[string]interface{}{
		"deal": *deal,
	}

	c.JSON(http.StatusOK, common.CreateSuccessResponse(mapObject))
}

type DealNetworkRequest struct {
	NetworkName string `json:"network_name"`
	DealId      string `json:"deal_id"`
}

type DealNetworkRequest_v1 struct {
	NetworkName string `json:"network_name"`
	DealId      int    `json:"deal_id"`
}
