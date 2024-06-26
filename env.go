package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func viperSetting() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetDefault("application.port", 8080)
	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		panic("讀取設定檔出現錯誤，原因為：" + viperErr.Error())
	}
}

func loadConfig() Configs {
	return Configs{
		dbUserName:   viper.GetString("db.username"),
		dbPassword:   viper.GetString("db.password"),
		dbConnection: viper.GetString("db.connection"),
	}
}

func connectDB(config Configs) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@%s",
		config.dbUserName, config.dbPassword, config.dbConnection)
	db, dbErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, dbErr
}

type Configs struct {
	dbUserName   string
	dbPassword   string
	dbConnection string
}
