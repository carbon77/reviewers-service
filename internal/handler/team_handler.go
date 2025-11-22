package handler

import (
	"errors"
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
	if errors.Is(err, errs.ResourceNotFound) {
		response := errs.NewErrorResponse("NOT_FOUND", err.Error())
		c.JSON(http.StatusNotFound, response)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body", "message": err.Error()})
		return
	}

	if err := h.service.CreateTeam(&team); err != nil {
		if errors.Is(err, errs.TeamAlreadyExists) {
			response := errs.NewErrorResponse("TEAM_EXISTS", fmt.Sprintf("%s already exists", team.Name))
			c.JSON(http.StatusBadRequest, response)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}
