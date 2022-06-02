package routers

import (
	"flink-data/common"
	"flink-data/common/constants"
	"flink-data/models"
	"flink-data/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Deal(router *gin.RouterGroup) {
	router.POST("", GetDeal)
	router.POST("latest", GetLatestDeal)
}

func GetDeal(c *gin.Context) {
	var dealNetworkRequest DealNetworkRequest
	err := c.ShouldBindBodyWith(&dealNetworkRequest, binding.JSON)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if dealNetworkRequest.DealId < 0 {
		err := fmt.Errorf("deal id:%d must be not less than 0", dealNetworkRequest.DealId)
		logs.GetLogger().Error(err)
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	networkName := dealNetworkRequest.NetworkName
	if networkName == "" {
		err := fmt.Errorf("network name must be provided")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if networkName != constants.NETWORK_CALIBRATION && networkName != constants.NETWORK_MAINNET {
		err := fmt.Errorf("network name should be %s or %s", constants.NETWORK_MAINNET, constants.NETWORK_CALIBRATION)
		logs.GetLogger().Error(err)
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	deal, err := service.GetDealById(dealNetworkRequest.DealId, networkName)
	if err != nil {
		if strings.Contains(err.Error(), "error code:-32603 message:pg: no rows in result set") {
			err := fmt.Errorf("deal not found")
			logs.GetLogger().Error(err)
			c.JSON(http.StatusNotFound, common.CreateErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		logs.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	if deal == nil {
		err := fmt.Errorf("deal not found")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusNotFound, common.CreateErrorResponse(http.StatusNotFound, err.Error()))
		return
	}

	mapObject := map[string]interface{}{
		"deal":   *deal,
		"result": deal.StoragePrice,
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
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	networkName := dealNetworkRequest.NetworkName
	if networkName == "" {
		err := fmt.Errorf("network id must be provided")
		logs.GetLogger().Error(err)
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if networkName != constants.NETWORK_CALIBRATION && networkName != constants.NETWORK_MAINNET {
		err := fmt.Errorf("network name should be %s or %s", constants.NETWORK_MAINNET, constants.NETWORK_CALIBRATION)
		logs.GetLogger().Error(err)
		c.JSON(http.StatusBadRequest, common.CreateErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	network, err := models.GetNetworkByName(networkName)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	if network == nil {
		err := fmt.Errorf("internal error, network:%s cannot be found", networkName)
		logs.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	deal, err := service.GetLatestDealByNetwork(network.Id)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	mapObject := map[string]interface{}{
		"deal": *deal,
	}

	c.JSON(http.StatusOK, common.CreateSuccessResponse(mapObject))
}

type DealNetworkRequest struct {
	NetworkName string `json:"network_name"`
	DealId      int64  `json:"deal_id"`
}
