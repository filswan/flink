package routers

import (
	"flink-data/common"
	"flink-data/common/constants"
	"flink-data/models"
	"fmt"
	"net/http"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
)

func Network(router *gin.RouterGroup) {
	router.GET(":network", GetNetwork)
}

func GetNetwork(c *gin.Context) {
	networkName := c.Param("network")

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
	c.JSON(http.StatusOK, common.CreateSuccessResponse(network))
}
