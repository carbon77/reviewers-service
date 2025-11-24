package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"reviewers/internal/errs"
	"reviewers/internal/models"
	"time"

	"gorm.io/gorm"
)

type PRRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewPRRepository(db *gorm.DB, logger *slog.Logger) *PRRepository {
	return &PRRepository{db, logger}
}

func (r *PRRepository) Create(pr *models.PullRequest) error {
	logger := r.logger.With(
		"method", "create_pull_request",
		"pull_request_id", pr.ID,
		"pull_request_name", pr.Name,
	)
	logger.Info("creating pull request")

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "Reviewers").Create(pr).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				logger.Warn("pull request already exists", "error", err)
				return errs.PullRequestExists
			} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
				logger.Warn("author not found", "error", err)
				return errs.ResourceNotFound
			}
			logger.Error("failed to create pull request", "error", err)
			return err
		}

		if err := tx.Model(pr).Association("Reviewers").Replace(pr.Reviewers); err != nil {
			return fmt.Errorf("failed to associate reviewers with pr %s: %w", pr.Name, err)
		}

		return nil
	})
}

func (r *PRRepository) Save(pr *models.PullRequest) error {
	logger := r.logger.With(
		"method", "update_pull_request",
		"pull_request_id", pr.ID,
		"pull_request_name", pr.Name,
	)
	logger.Info("updating pull request")

	err := r.db.Save(&pr).Error
	if err != nil {
		logger.Error("failed to update pull request", "error", err)
		return err
	}

	if err := r.db.Model(pr).Association("Reviewers").Replace(pr.Reviewers); err != nil {
		logger.Error("failed to update pull request", "error", err)
		return fmt.Errorf("failed to associate reviewers with pr %s: %w", pr.Name, err)
	}

	return nil
}

func (r *PRRepository) Get(pullRequestID string) (*models.PullRequest, error) {
	logger := r.logger.With(
		"method", "get_pull_request",
		"pull_request_id", pullRequestID,
	)
	logger.Info("getting pull request")

	var pr models.PullRequest
	err := r.db.Where("pull_request_id = ?", pullRequestID).Preload("Reviewers").First(&pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("pull request not found", "error", err)
			return nil, errs.ResourceNotFound
		}
		logger.Error("failed to get pull request", "error", err)
		return nil, err
	}

	return &pr, nil
}

func (r *PRRepository) Merge(pullRequestID string) error {
	logger := r.logger.With(
		"method", "merge_pull_request",
		"pull_request_id", pullRequestID,
	)
	logger.Info("merging pull request")

	now := time.Now()
	pr := models.PullRequest{
		ID:       pullRequestID,
		Status:   models.StatusMerged,
		MergedAt: &now,
	}

	err := r.db.Model(&pr).Updates(pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("pull request not found", "error", err)
			return errs.ResourceNotFound
		}
		logger.Error("failed to merge pull request", "error", err)
		return err
	}

	return nil
}
