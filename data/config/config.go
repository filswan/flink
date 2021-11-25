package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/filswan/go-swan-lib/logs"
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

func InitConfig() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		logs.GetLogger().Fatal("Cannot get home directory.")
	}

	configFile := filepath.Join(homedir, ".swan/filink/data/config.toml")

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
