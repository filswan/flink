package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/filswan/go-swan-lib/logs"
)

type Configuration struct {
	Port      int       `toml:"port"`
	Release   bool      `toml:"release"`
	Database  database  `toml:"database"`
	ChainLink chainLink `toml:"chain_link"`
}

type database struct {
	DbHost           string `toml:"db_host"`
	DbPort           int    `toml:"db_port"`
	DbSchemaName     string `toml:"db_schema_name"`
	DbUsername       string `toml:"db_username"`
	DbPassword       string `toml:"db_password"`
	DbArgs           string `toml:"db_args"`
	DbMaxIdleConnNum int    `toml:"db_max_idle_conn_num"`
}

type chainLink struct {
	BulkInsertChainlinkLimit   int   `toml:"bulk_insert_chainlink_limit"`
	BulkInsertIntervalMilliSec int64 `toml:"bulk_insert_interval_milli_sec"`
	DealIdIntervalMax          int64 `toml:"deal_id_interval_max"`
}

var config *Configuration

func InitConfig() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		logs.GetLogger().Fatal("Cannot get home directory.")
	}

	configFile := filepath.Join(homedir, ".swan/flink/data/config.toml")
	//configFile := filepath.Join(homedir, "Documents/NebulaAI/Go-Tutorial/flink/data/config/config.toml")

	metaData, err := toml.DecodeFile(configFile, &config)

	if err != nil {
		log.Fatal("error:", err)
	}

	if !requiredFieldsAreGiven(metaData) {
		log.Fatal("required fields not given")
	}
}

func GetConfig() Configuration {
	if config == nil {
		InitConfig()
	}
	return *config
}

func requiredFieldsAreGiven(metaData toml.MetaData) bool {
	requiredFields := [][]string{
		{"port"},
		{"release"},
		{"database"},
		{"chain_link"},

		{"database", "db_host"},
		{"database", "db_port"},
		{"database", "db_schema_name"},
		{"database", "db_username"},
		{"database", "db_password"},
		{"database", "db_args"},
		{"database", "db_max_idle_conn_num"},

		{"chain_link", "bulk_insert_chainlink_limit"},
		{"chain_link", "bulk_insert_interval_milli_sec"},
		{"chain_link", "deal_id_interval_max"},
	}

	for _, v := range requiredFields {
		if !metaData.IsDefined(v...) {
			log.Fatal("required fields ", v)
		}
	}

	return true
}
