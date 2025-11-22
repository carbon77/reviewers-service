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

	// Users
	userRepository := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	userRouter := router.Group("/users")
	userRouter.POST("/setIsActive", userHandler.SetActiveStatus)
	userRouter.GET("/getReview", userHandler.GetReview)

	// Teams
	teamRepository := repository.NewTeamRepository(conn)
	teamService := service.NewTeamService(teamRepository)
	teamHandler := handler.NewTeamHandler(teamService)

	teamRouter := router.Group("/team")
	teamRouter.GET("/get", teamHandler.GetTeam)
	teamRouter.POST("/add", teamHandler.CreateTeam)

	// Pull requests
	prRepository := repository.NewPRRepository(conn)
	prService := service.NewPRService(prRepository, teamService)
	prHandler := handler.NewPRHandler(prService)

	prRouter := router.Group("/pullRequest")
	prRouter.POST("/create", prHandler.Create)

	addr := fmt.Sprintf(":%d", cfg.Port)
	router.Run(addr)
}
