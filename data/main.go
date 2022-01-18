package main

import (
	"flink-data/common/constants"
	"flink-data/config"
	"flink-data/database"
	"flink-data/routers"
	"flink-data/service"
	"strconv"
	"time"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	/*
		if len(os.Args) < 2 {
			logs.GetLogger().Fatal("Flink network must be specified (calibration/mainnet) ")
		}*/

	db := database.Init()
	defer database.CloseDB(db)

	subCmd := constants.PARAM_CALIBRATION // os.Args[1]
	switch subCmd {
	case constants.PARAM_CALIBRATION:
		logs.GetLogger().Info("starting for calibration network")
		go service.GetDealsFromCalibrationLoop()
		createGinServer()
	case constants.PARAM_MAINNET:
		logs.GetLogger().Info("starting for mainnet network")
		go service.GetDealsFromMainnet()
		createGinServer()
	default:
		logs.GetLogger().Fatal("sub command should be: calibration|mainnet")
	}
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
