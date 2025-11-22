package main

import (
	"fmt"
	"reviewers/internal/config"
	"reviewers/internal/db"
	"reviewers/internal/handler"
	"reviewers/internal/repository"
	"reviewers/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	conn := db.Connect(cfg)

	router := gin.Default()

	{
		userRepository := repository.NewUserRepository(conn)
		userService := service.NewUserService(userRepository)
		userHandler := handler.NewUserHandler(userService)

		userRouter := router.Group("/users")
		userRouter.POST("/setIsActive", userHandler.SetActiveStatus)
		userRouter.GET("/getReview", userHandler.GetReview)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	router.Run(addr)
}
