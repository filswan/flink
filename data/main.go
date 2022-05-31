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
	if len(os.Args) < 2 {
		logs.GetLogger().Fatal("Flink network must be specified")
	}

	network := os.Args[1]
	if network != constants.PARAM_CALIBRATION && network != constants.PARAM_MAINNET {
		err := fmt.Errorf("network should be: %s|%s", constants.PARAM_CALIBRATION, constants.PARAM_MAINNET)
		logs.GetLogger().Fatal(err)
	}

	logs.GetLogger().Info("starting for ", network, " network")
	setConfigFilepath(network)

	db := database.Init()
	defer database.CloseDB(db)

	if network == constants.PARAM_CALIBRATION {
		go service.GetDealsFromCalibrationLoop()
	} else {
		go service.GetDealsFromMainnetLoop()
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

	if configFilepath == nil || *configFilepath == "" {
		logs.GetLogger().Info("you do not provide config file, use default")
	} else {
		logs.GetLogger().Info("config file you provided is:", *configFilepath)
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
