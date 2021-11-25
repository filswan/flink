package database

import (
	"filecoin-data-provider/data/config"
	"strconv"

	"github.com/filswan/go-swan-lib/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	confDb := config.GetConfig().Database
	dbSource := confDb.DbUsername + ":" + confDb.DbPassword + "@tcp(" + confDb.DbHost + ":" + strconv.Itoa(confDb.DbPort) + ")/" + confDb.DbSchemaName + "?" + confDb.DbArgs
	db, err := gorm.Open("mysql", dbSource)
	if err != nil {
		logs.GetLogger().Fatal("db err: ", err)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(confDb.DbMaxIdleConnNum)
	//db.LogMode(true)
	db.LogMode(false)
	DB = db

	return DB
}

// Using this function to get a connection
// You can create your connection pool here.
func GetDB() *gorm.DB {
	if DB == nil {
		DB = Init()
	}

	return DB
}

func SaveOne(data interface{}) error {
	db := GetDB()
	err := db.Save(data).Error
	return err
}

func CloseDB(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		logs.GetLogger().Error(err)
	}
}
