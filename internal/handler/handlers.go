package handler

import (
	"log/slog"
	"reviewers/internal/repository"
	"reviewers/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitHandlers(logger *slog.Logger, conn *gorm.DB, router *gin.Engine) {
	// Users
	userRepository := repository.NewUserRepository(conn, logger)
	userService := service.NewUserService(userRepository)
	userHandler := NewUserHandler(userService)

	userRouter := router.Group("/users")
	userRouter.POST("/setIsActive", userHandler.SetActiveStatus)
	userRouter.GET("/getReview", userHandler.GetReview)

	// Teams
	teamRepository := repository.NewTeamRepository(conn, logger)
	teamService := service.NewTeamService(teamRepository)
	teamHandler := NewTeamHandler(teamService)

	teamRouter := router.Group("/team")
	teamRouter.GET("/get", teamHandler.GetTeam)
	teamRouter.POST("/add", teamHandler.CreateTeam)
	teamRouter.POST("/deactivate", teamHandler.DeactivateTeam)

	// Pull requests
	prRepository := repository.NewPRRepository(conn, logger)
	prService := service.NewPRService(prRepository, teamService, userService)
	prHandler := NewPRHandler(prService)

	prRouter := router.Group("/pullRequest")
	prRouter.POST("/create", prHandler.Create)
	prRouter.POST("/merge", prHandler.Merge)
	prRouter.POST("/reassign", prHandler.Reassign)

}
