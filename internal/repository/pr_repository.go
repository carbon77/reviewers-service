package repository

import (
	"errors"
	"fmt"
	"reviewers/internal/errs"
	"reviewers/internal/models"
	"time"

	"gorm.io/gorm"
)

type PRRepository struct {
	db *gorm.DB
}

func NewPRRepository(db *gorm.DB) *PRRepository {
	return &PRRepository{db}
}

func (r *PRRepository) Create(pr *models.PullRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "Reviewers").Create(pr).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errs.PullRequestExists
			} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return errs.ResourceNotFound
			}
			return err
		}

		if err := tx.Model(pr).Association("Reviewers").Replace(pr.Reviewers); err != nil {
			return fmt.Errorf("failed to associate reviewers with pr %s: %w", pr.Name, err)
		}

		return nil
	})
}

func (r *PRRepository) Save(pr *models.PullRequest) error {
	err := r.db.Save(&pr).Error
	if err != nil {
		return err
	}

	if err := r.db.Model(pr).Association("Reviewers").Replace(pr.Reviewers); err != nil {
		return fmt.Errorf("failed to associate reviewers with pr %s: %w", pr.Name, err)
	}

	return nil
}

func (r *PRRepository) Get(pullRequestID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := r.db.Where("pull_request_id = ?", pullRequestID).Preload("Reviewers").First(&pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ResourceNotFound
		}
		return nil, err
	}

	for _, user := range pr.Reviewers {
		pr.AssignedReviewers = append(pr.AssignedReviewers, user.ID)
	}

	return &pr, nil
}

func (r *PRRepository) Merge(pullRequestID string) error {
	now := time.Now()
	pr := models.PullRequest{
		ID:       pullRequestID,
		Status:   models.StatusMerged,
		MergedAt: &now,
	}

	err := r.db.Model(&pr).Updates(pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ResourceNotFound
		}
		return err
	}

	return nil
}
