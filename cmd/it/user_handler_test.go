package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reviewers/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSetActiveStatus(t *testing.T) {
	runInTransaction(t, func(tx *gorm.DB) {
		r := setupRouter(tx)
		createTestUser(t, tx)

		reqBody, _ := json.Marshal(map[string]interface{}{
			"user_id":   user.ID,
			"is_active": false,
		})

		req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "status updated", resp["message"])

		var updatedUser models.User
		tx.First(&updatedUser, "user_id = ?", user.ID)
		assert.False(t, updatedUser.IsActive)
	})
}

func TestSetActiveStatus_UserNotFound(t *testing.T) {
	runInTransaction(t, func(tx *gorm.DB) {
		r := setupRouter(tx)

		reqBody, _ := json.Marshal(map[string]interface{}{
			"user_id":   uuid.New().String(),
			"is_active": false,
		})

		req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		errorResponse := resp["error"].(map[string]interface{})
		assert.Equal(t, "NOT_FOUND", errorResponse["code"])
		assert.Equal(t, "resource not found", errorResponse["message"])
	})
}

func TestSetActiveStatus_BadRequest(t *testing.T) {
	runInTransaction(t, func(tx *gorm.DB) {
		r := setupRouter(tx)

		reqBody, _ := json.Marshal(map[string]interface{}{})

		req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "invalid body", resp["error"])
	})
}

func TestGetReview(t *testing.T) {
	runInTransaction(t, func(tx *gorm.DB) {
		r := setupRouter(tx)
		createTestUser(t, tx)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/users/getReview?user_id=%s", user.ID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, user.ID, resp["user_id"])
		assert.Len(t, resp["pull_requests"], 1)

		prs := resp["pull_requests"].([]interface{})
		assert.Equal(t, map[string]interface{}{
			"author_id":         pr.AuthorID,
			"pull_request_id":   pr.ID,
			"pull_request_name": pr.Name,
			"status":            pr.Status,
		}, prs[0])
	})
}

var user models.User
var team models.Team
var pr models.PullRequest

func createTestUser(t *testing.T, tx *gorm.DB) {
	teamID := uuid.New().String()
	team = models.Team{
		ID:   teamID,
		Name: "test_team",
	}
	tx.Create(&team)

	userID := uuid.New().String()
	user = models.User{
		ID:       userID,
		Username: "test",
		IsActive: true,
		TeamID:   teamID,
	}
	tx.Create(&user)

	prID := "pr-0001"
	reviewer := models.PullRequestReviewer{PullRequestID: prID, UserID: userID}
	pr = models.PullRequest{
		ID:        prID,
		Name:      "test",
		Status:    models.StatusOpen,
		Reviewers: []models.PullRequestReviewer{reviewer},
		Author:    user,
	}
	tx.Create(&pr)
}
