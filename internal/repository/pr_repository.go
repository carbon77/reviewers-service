package repository

import (
	"errors"
	"fmt"
	"reviewers/internal/errs"
	"reviewers/internal/models"

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
