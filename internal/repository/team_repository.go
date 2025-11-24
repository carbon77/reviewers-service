package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"reviewers/internal/errs"
	"reviewers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTeamRepository(db *gorm.DB, logger *slog.Logger) *TeamRepository {
	return &TeamRepository{db, logger}
}

func (r *TeamRepository) GetTeam(name string) (*models.Team, error) {
	logger := r.logger.With(
		"method", "get_team",
		"team_name", name,
	)
	logger.Info("getting team")

	var team models.Team

	err := r.db.Where("name = ?", name).Preload("Members").First(&team).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warn("team not found", "error", err)
		return &team, errs.ResourceNotFound
	}

	return &team, err
}

func (r *TeamRepository) CreateTeam(team *models.Team) error {
	logger := r.logger.With(
		"method", "create_team",
		"team_name", team.Name,
	)
	logger.Info("creating team")

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create team
		team.ID = uuid.New().String()
		if err := tx.Select("ID", "Name").Create(team).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				logger.Warn("team already exists", "error", err)
				return errs.TeamExists
			}
			logger.Error("failed to create team", "error", err)
			return fmt.Errorf("failed to create team %s: %w", team.Name, err)
		}

		// Create users
		var newUsers []*models.User
		for i := range team.Members {
			user := &team.Members[i]

			if user.ID != "" {
				if err := tx.Save(user).Error; err != nil {
					logger.Error("failed to update user", "error", err, "user_id", user.ID)
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
				logger.Error("failed to create users", "error", err, "user_id")
				return fmt.Errorf("failed to create users: %w", err)
			}
		}

		return nil
	})
}

func (r *TeamRepository) GetReviewerIdsFromUserTeam(userID string, excludedUsers ...string) ([]*models.User, error) {
	logger := r.logger.With(
		"method", "get_reviewers_from_same_team",
		"user_id", userID,
	)
	logger.Info("getting reviewers from the same team")

	var reviewers []*models.User

	excludedIds := make([]string, 0, len(excludedUsers)+1)
	excludedIds = append(excludedIds, excludedUsers...)
	excludedIds = append(excludedIds, userID)

	query := r.db.Model(&models.User{}).
		Where("team_id = (?)", r.db.Model(&models.User{}).
			Select("team_id").
			Where("user_id = ?", userID).
			Limit(1),
		).
		Where("is_active = true").
		Where("user_id NOT IN ?", excludedIds)

	err := query.Find(&reviewers).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("team not found", "error", err)
			return nil, errs.ResourceNotFound
		}
		logger.Error("failed to find users from team", "error", err)
		return nil, fmt.Errorf("failed to find users from team: %s", err.Error())
	}

	return reviewers, nil
}

func (r *TeamRepository) DeactivateTeam(teamID string) error {
	logger := r.logger.With(
		"method", "deactivate_team",
		"team_id", teamID,
	)
	logger.Info("deactivating team")

	err := r.db.Model(&models.User{}).Where("team_id = ?", teamID).Update("is_active", false).Error
	if err != nil {
		logger.Error("failed to deactivate team", "error", err)
	}

	return err
}
