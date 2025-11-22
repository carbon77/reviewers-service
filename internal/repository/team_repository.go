package repository

import (
	"errors"
	"fmt"
	"reviewers/internal/errs"
	"reviewers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db}
}

func (r *TeamRepository) GetTeam(name string) (*models.Team, error) {
	var team models.Team

	result := r.db.Where("name = ?", name).Preload("Members").First(&team)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &team, errs.ResourceNotFound
	}

	return &team, result.Error
}

func (r *TeamRepository) CreateTeam(team *models.Team) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create team
		team.ID = uuid.New().String()
		if err := tx.Select("ID", "Name").Create(team).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errs.TeamAlreadyExists
			}
			return fmt.Errorf("failed to create team %s: %w", team.Name, err)
		}

		// Create users
		var newUsers []*models.User
		for i := range team.Members {
			user := &team.Members[i]

			if user.ID != "" {
				if err := tx.Save(user).Error; err != nil {
					return fmt.Errorf("failed to update user %s: %w", user.ID, err)
				}
			} else {
				user.ID = uuid.New().String()
				user.TeamID = team.ID
				newUsers = append(newUsers, user)
			}
		}

		if len(newUsers) > 0 {
			if err := tx.Create(newUsers).Error; err != nil {
				return fmt.Errorf("failed to create users: %w", err)
			}
		}

		return nil
	})
}

func (r *TeamRepository) GetReviewerIdsFromUserTeam(userID string) ([]*models.User, error) {
	var reviewers []*models.User

	err := r.db.Model(&models.User{}).
		Where("team_id = (?)", r.db.Model(&models.User{}).
			Select("team_id").
			Where("user_id = ?", userID).
			Limit(1),
		).
		Where("user_id <> ?", userID).
		Where("is_active = true").
		Find(&reviewers).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ResourceNotFound
		}
		return nil, fmt.Errorf("failed to find users from team: %s", err.Error())
	}

	return reviewers, nil
}
