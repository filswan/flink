package routers

import (
	"filecoin-data-provider/data/common"
	"filecoin-data-provider/data/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Common(router *gin.RouterGroup) {
	router.GET("host", GetHostInfo)
}

func GetHostInfo(c *gin.Context) {
	info := service.GetHostInfo()
	c.JSON(http.StatusOK, common.CreateSuccessResponse(info))
}
