package main

import (
	"flag"
	"flink-data/common/constants"
	"flink-data/config"
	"flink-data/database"
	"flink-data/routers"
	"flink-data/service"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	network := constants.PARAM_CALIBRATION //default

	if len(os.Args) >= 2 {
		network = os.Args[1]
	}

	logs.GetLogger().Info("starting get deals network:", network)
	setConfigFilepath(network)

	db := database.Init()
	defer database.CloseDB(db)

	switch network {
	case constants.PARAM_MAINNET:
		go service.GetDealsFromMainnetLoop()
	default:
		go service.GetDealsFromCalibrationLoop()
	}

	createGinServer()
}

func setConfigFilepath(subCmdName string) error {
	cmd := flag.NewFlagSet(subCmdName, flag.ExitOnError)

	configFilepath := cmd.String("c", "", "config file path")

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	if !cmd.Parsed() {
		err := fmt.Errorf("sub command parse failed")
		logs.GetLogger().Error(err)
		return err
	}

	config.InitConfig(configFilepath)

	return nil
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
	routers.Network(ginEngine.Group("/network"))
	err := ginEngine.Run(":" + strconv.Itoa(config.GetConfig().Port))
	if err != nil {
		logs.GetLogger().Fatal(err)
	}
}
