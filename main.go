package main

import (
	categoryHttp "GoProject/module/category/delivery"
	categoryRepository "GoProject/module/category/repository"
	category "GoProject/module/category/service"
	loginHttp "GoProject/module/login/delivery"
	loginSvc "GoProject/module/login/service"
	userHttp "GoProject/module/user/delivery"
	userRepository "GoProject/module/user/repository"
	userSvc "GoProject/module/user/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	// CORS configuration
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	//gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// Apply the CORS middleware to the router
	server.Use(cors.New(corsConfig))

	userRepo := userRepository.NewUserRepository(db)
	categoryRepo := categoryRepository.NewCategoryRepository(db)
	categorySvc := category.NewCategoryService(categoryRepo)
	userService := userSvc.NewUserService(userRepo, categorySvc)
	loginService := loginSvc.NewLoginService(userService)
	userHttp.NewUserHttpHandler(userService, server)
	loginHttp.NewLoginHttpHandler(loginService, server)
	categoryHttp.NewCategoryHandler(categorySvc, userService, server)

	server.GET("/go", func(context *gin.Context) {
		var intChen = make(chan int, 3)
		go aync(intChen)
		go aync2(intChen)
	})

	err := server.Run(":" + viper.GetString("application.port"))
	if err != nil {
		log.Println("Error Occur!")
		return
	}

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
