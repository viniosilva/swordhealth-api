package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viniosilva/swordhealth-api/docs"
	"github.com/viniosilva/swordhealth-api/internal/config"
	"github.com/viniosilva/swordhealth-api/internal/controller"
	"github.com/viniosilva/swordhealth-api/internal/repository"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

// @title           Sword Health API
// @version         1.0
// @description     Task manager test
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
// @BasePath  /api
func main() {
	log.SetFormatter(&log.JSONFormatter{})

	c := config.LoadConfig()
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
		c.MySQL.Username, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.Database))
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	router := r.Group("/api")

	healthRepository := repository.NewHealthRepository(db)
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)

	cryptoService := service.NewCryptoService(c.Crypto.HashKey, c.Crypto.JwtKey, c.Crypto.ExpiresIn)
	healthService := service.NewHealthService(healthRepository)
	userService := service.NewUserService(userRepository)
	taskService := service.NewTaskService(taskRepository)
	notificationService := service.NewNotificationService(userRepository)
	authService := service.NewAuthService()

	controller.NewHealthController(router, healthService)
	controller.NewUserController(router, userService, cryptoService)
	controller.NewTaskController(router, taskService, notificationService)
	controller.NewAuthController(router, authService, userService, cryptoService)

	host := fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
	docs.SwaggerInfo.Host = host
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(host)
}
