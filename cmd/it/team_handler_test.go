package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reviewers/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTeam(t *testing.T) {
	runInTransaction(t, func(tx *gorm.DB) {
		r := setupRouter(tx)
		team := models.Team{
			ID:   uuid.New().String(),
			Name: "test1",
			Members: []models.User{
				{
					ID:       uuid.New().String(),
					Username: "dima",
					IsActive: false,
				},
			},
		}
		tx.Create(&team)

		teamBody := map[string]interface{}{
			"team_name": "test",
			"members": []map[string]interface{}{
				{
					"username":  "igor",
					"is_active": true,
				},
				{
					"user_id":   team.Members[0].ID,
					"is_active": true,
				},
			},
		}
		reqBody, _ := json.Marshal(teamBody)
		req, _ := http.NewRequest("POST", "/team/add", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Team exists
		req, _ = http.NewRequest("POST", "/team/add", bytes.NewBuffer(reqBody))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, map[string]interface{}{
			"code":    "TEAM_EXISTS",
			"message": "test already exists",
		}, resp["error"])

		// Check team created
		var createdTeam models.Team
		tx.Select("Name").First(&createdTeam).Preload("Members")
		assert.Equal(t, teamBody["team_name"], createdTeam.Name)

		var user1 models.User
		tx.Where("user_id = ?", team.Members[0].ID).First(&user1)
		assert.NotEqual(t, team.ID, user1.TeamID)
		assert.Equal(t, true, user1.IsActive)

		var user2 models.User
		tx.Where("user_id <> ?", team.Members[0].ID).First(&user2)
		assert.NotNil(t, user2)
		assert.Equal(t, "igor", user2.Username)
		assert.Equal(t, true, user2.IsActive)
		assert.Equal(t, user1.TeamID, user2.TeamID)
	})
}

func TestGetTeam(t *testing.T) {
	runInTransaction(t, func(tx *gorm.DB) {
		r := setupRouter(tx)

		team := models.Team{
			ID:   uuid.New().String(),
			Name: "test",
			Members: []models.User{
				{
					ID:       uuid.New().String(),
					Username: "igor",
					IsActive: true,
				},
			},
		}
		tx.Create(&team)

		// Success
		req, _ := http.NewRequest("GET", "/team/get?team_name=test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, map[string]interface{}{
			"team_name": "test",
			"members": []interface{}{
				map[string]interface{}{
					"user_id":   team.Members[0].ID,
					"username":  team.Members[0].Username,
					"is_active": team.Members[0].IsActive,
				},
			},
		}, resp)

		// Not found
		req, _ = http.NewRequest("GET", "/team/get?team_name=123", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, map[string]interface{}{
			"code":    "NOT_FOUND",
			"message": "resource not found",
		}, resp["error"])
	})
}
