package repository

import (
	"reviewers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) SetActiveStatus(userID uuid.UUID, active bool) error {
	return r.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("is_active", active).Error
}

func (r *UserRepository) GetReview(userID uuid.UUID) ([]models.PullRequest, error) {
	var prs []models.PullRequest

	err := r.db.Model(&models.PullRequest{}).
		Joins("JOIN pull_request_reviewers prr ON prr.pull_request_id = pull_requests.pull_request_id").
		Where("prr.user_id = ?", userID).
		Find(&prs).Error

	return prs, err
}
