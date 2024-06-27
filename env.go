package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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

func loadServer() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	server.Use(cors.New(corsConfig))
	return server
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
