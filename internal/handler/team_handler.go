package handler

import (
	"fmt"
	"net/http"
	"reviewers/internal/errs"
	"reviewers/internal/models"
	"reviewers/internal/service"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service *service.TeamService
}

func NewTeamHandler(service *service.TeamService) *TeamHandler {
	return &TeamHandler{service}
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	name := c.Query("team_name")

	team, err := h.service.GetTeam(name)
	if err != nil {
		switch err.(type) {
		case errs.ApiError:
			err.(errs.ApiError).ReturnError(c, err.Error())
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something bad happened"})
		}
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.service.CreateTeam(&team); err != nil {
		switch err.(type) {
		case errs.ApiError:
			err.(errs.ApiError).ReturnError(c, fmt.Sprintf("%s already exists", team.Name))
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something bad happened"})
		}
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) DeactivateTeam(c *gin.Context) {
	teamID := c.Query("team_id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	err := h.service.DeactivateTeam(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something bad happened"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "team has been deactivated"})
}
