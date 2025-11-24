package repository

import (
	"errors"
	"reviewers/internal/errs"
	"reviewers/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) SetActiveStatus(userID string, active bool) error {
	result := r.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("is_active", active)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ResourceNotFound
	}

	return nil
}

func (r *UserRepository) GetReview(userID string) ([]models.PullRequestShort, error) {
	var prs []models.PullRequestShort

	err := r.db.Model(&models.PullRequest{}).
		Select("pull_requests.pull_request_id", "Name", "AuthorID", "Status").
		Joins("JOIN pull_request_reviewers prr ON prr.pull_request_id = pull_requests.pull_request_id").
		Where("prr.user_id = ?", userID).
		Find(&prs).Error

	return prs, err
}

func (r *UserRepository) Get(userID string) (*models.User, error) {
	var user models.User

	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ResourceNotFound
		}
		return nil, err
	}

	return &user, nil
}
