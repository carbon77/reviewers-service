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

type PRHandler struct {
	service *service.PRService
}

func NewPRHandler(service *service.PRService) *PRHandler {
	return &PRHandler{service}
}

type CreatePRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func (h *PRHandler) Create(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	pr := &models.PullRequest{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}

	if err := h.service.Create(pr); err != nil {
		if errors.Is(err, errs.ResourceNotFound) {
			response := errs.NewErrorResponse("NOT_FOUND", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		} else if errors.Is(err, errs.PullRequestExists) {
			response := errs.NewErrorResponse("PR_EXISTS", fmt.Sprintf("PR %s already exists", pr.ID))
			c.JSON(http.StatusConflict, response)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}
