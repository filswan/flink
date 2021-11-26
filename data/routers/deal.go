package routers

import (
	"filink/data/common"
	"filink/data/service"
	"net/http"
	"strconv"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
)

func Deal(router *gin.RouterGroup) {
	router.GET(":deal_id", GetDeal)
}

func GetDeal(c *gin.Context) {
	dealIdStr := c.Param("deal_id")
	dealId, err := strconv.ParseInt(dealIdStr, 10, 64)
	if err != nil {
		logs.GetLogger().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse(err.Error()))
		return
	}

	deal, err := service.GetDealById(dealId)
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

func AddMiner(c *gin.Context) {
}
