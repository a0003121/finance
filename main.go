package main

import (
	"GoProject/handler"
	categoryRepository "GoProject/module/category/repository"
	category "GoProject/module/category/service"
	loginSvc "GoProject/module/login/service"
	userRepository "GoProject/module/user/repository"
	userSvc "GoProject/module/user/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	viperSetting()
	config := loadConfig()
	db, dbErr := connectDB(config)
	if dbErr != nil {
		panic("connect to db error: " + dbErr.Error())
	}

	// Auto migrate to keep schema updated
	migrateErr := db.AutoMigrate()
	if migrateErr != nil {
		panic("auto migrate error: " + migrateErr.Error())
		return
	}

	server := loadServer()
	registerHandler(db, server)

	err := server.Run(":" + viper.GetString("application.port"))
	if err != nil {
		log.Println("Error Occur!")
		return
	}

}

func registerHandler(db *gorm.DB, server *gin.Engine) {
	userRepo := userRepository.NewUserRepository(db)
	categoryRepo := categoryRepository.NewCategoryRepository(db)
	categorySvc := category.NewCategoryService(categoryRepo)
	userService := userSvc.NewUserService(userRepo, categorySvc)
	loginService := loginSvc.NewLoginService(userService)
	handler.NewUserHttpHandler(userService, server)
	handler.NewLoginHttpHandler(loginService, server)
	handler.NewCategoryHandler(categorySvc, userService, server)
	handler.NewRecordHandler(categorySvc, userService, server)
	handler.NewExcelHttpHandler(userService, categorySvc, server)
	handler.NewStatisticsHttpHandler(categorySvc, userService, server)

	server.GET("/go", func(context *gin.Context) {
		var intChen = make(chan int, 3)
		go aync(intChen)
		go aync2(intChen)
	})
}

func aync(intChen chan int) {
	for i := 0; i < 5; i++ {
		intChen <- i
		log.Printf("gogo! i  %d", i)
		time.Sleep(time.Second)
	}
}

func aync2(intChen chan int) {
	for i := 0; i < 5; i++ {
		log.Printf("gogo! chan  %d", <-intChen)
		time.Sleep(time.Millisecond)
	}
}
