package main

import (
	"filecoin-data-provider/data/config"
	"filecoin-data-provider/data/database"
	"filecoin-data-provider/data/routers"
	"filecoin-data-provider/data/service"
	"strconv"
	"time"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	//test.Test()

	db := database.Init()
	defer database.CloseDB(db)

	go service.GetDealsFromCalibrationLoop()

	createGinServer()
}

func createGinServer() {
	if config.GetConfig().Release {
		gin.SetMode(gin.ReleaseMode)
	}

	ginEngine := gin.Default()
	ginEngine.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	routers.Common(ginEngine.Group("/common"))
	routers.Deal(ginEngine.Group("/deal"))
	err := ginEngine.Run(":" + strconv.Itoa(config.GetConfig().Port))
	if err != nil {
		logs.GetLogger().Fatal(err)
	}
}
