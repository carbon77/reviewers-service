package repository

import (
	"errors"
	"log/slog"
	"reviewers/internal/errs"
	"reviewers/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{db, logger}
}

func (r *UserRepository) SetActiveStatus(userID string, active bool) error {
	logger := r.logger.With(
		"method", "set_active_status",
		"user_id", userID,
	)
	logger.Info("setting active status")

	result := r.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("is_active", active)

	if result.Error != nil {
		logger.Error("failed to set status", "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Warn("user not found", "error", result.Error)
		return errs.ResourceNotFound
	}

	return nil
}

func (r *UserRepository) GetReview(userID string) ([]models.PullRequestShort, error) {
	logger := r.logger.With(
		"method", "get_reviews",
		"user_id", userID,
	)
	logger.Info("getting reviews")

	var prs []models.PullRequestShort

	err := r.db.Model(&models.PullRequest{}).
		Select("pull_requests.pull_request_id", "Name", "AuthorID", "Status").
		Joins("JOIN pull_request_reviewers prr ON prr.pull_request_id = pull_requests.pull_request_id").
		Where("prr.user_id = ?", userID).
		Find(&prs).Error
	if err != nil {
		logger.Error("failed to get reviews", "error", err)
	}

	return prs, err
}

func (r *UserRepository) Get(userID string) (*models.User, error) {
	logger := r.logger.With(
		"method", "get_user",
		"user_id", userID,
	)
	logger.Info("getting user")

	var user models.User

	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("user not found", "error", err)
			return nil, errs.ResourceNotFound
		}
		logger.Error("failed to get user", "error", err)
		return nil, err
	}

	return &user, nil
}
