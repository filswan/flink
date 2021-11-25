package config

import (
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Port     int      `toml:"port"`
	Release  bool     `toml:"release"`
	Database database `toml:"database"`
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

var config *Configuration

func InitConfig(configFile string) {
	if strings.Trim(configFile, " ") == "" {
		configFile = "./config/config.toml"
	}
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
		InitConfig("")
	}
	return *config
}

func GetConfigFromMainParams(configFile string) Configuration {
	if config == nil {
		InitConfig(configFile)
	}
	return *config
}

func requiredFieldsAreGiven(metaData toml.MetaData) bool {
	requiredFields := [][]string{
		{"port"},
		{"release"},
		{"database"},

		{"database", "db_host"},
		{"database", "db_port"},
		{"database", "db_schema_name"},
		{"database", "db_username"},
		{"database", "db_password"},
		{"database", "db_args"},
		{"database", "db_max_idle_conn_num"},
	}

	for _, v := range requiredFields {
		if !metaData.IsDefined(v...) {
			log.Fatal("required fields ", v)
		}
	}

	return true
}
