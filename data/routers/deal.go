package routers

import (
	"flink-data/common"
	"flink-data/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
)

func Deal(router *gin.RouterGroup) {
	router.POST("", GetDeal)
	router.GET("latest", GetLatestDeal)
}

func GetDeal(c *gin.Context) {
	var dealNetworkRequest DealNetworkRequest
	err := c.BindJSON(&dealNetworkRequest)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
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
	networkIdStr := c.Query("network_id")
	if networkIdStr == "" {
		err := fmt.Errorf("network id must be provided")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	networkId, err := strconv.ParseInt(networkIdStr, 10, 64)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	deal, err := service.GetLatestDealByNetwork(networkId)
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
