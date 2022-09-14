package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/viniosilva/swordhealth-api/internal/config"
	"github.com/viniosilva/swordhealth-api/internal/controller"
	"github.com/viniosilva/swordhealth-api/internal/repository"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

func main() {
	c := config.LoadConfig()
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s",
		c.MySQL.Username, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.Database))
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	router := r.Group("/api")

	healthRepository := repository.NewHealthRepository(db)
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)

	healthService := service.NewHealthService(healthRepository)
	userService := service.NewUserService(userRepository)
	taskService := service.NewTaskService(taskRepository)

	controller.NewHealthController(router, healthService)
	controller.NewUserController(router, userService)
	controller.NewTaskController(router, taskService)

	r.Run(fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port))
}
