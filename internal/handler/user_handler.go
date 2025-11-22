package handler

import (
	"errors"
	"net/http"
	"reviewers/internal/errs"
	"reviewers/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

type SetActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) SetActiveStatus(c *gin.Context) {
	var req SetActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	userUuid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.service.SetActiveStatus(userUuid, req.IsActive); err != nil {
		if errors.Is(err, errs.ResourceNotFound) {
			response := errs.NewErrorResponse("NOT_FOUND", err.Error())
			c.JSON(http.StatusNotFound, response)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated",
	})
}

func (h *UserHandler) GetReview(c *gin.Context) {
	userId := c.Query("user_id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	userUuid, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	prs, err := h.service.GetReview(userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userId,
		"pull_requests": prs,
	})
}
