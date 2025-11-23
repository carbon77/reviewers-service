package service

import (
	"reviewers/internal/models"
	"reviewers/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) SetActiveStatus(userID string, active bool) error {
	return s.repo.SetActiveStatus(userID, active)
}

func (s *UserService) GetReview(userID string) ([]models.PullRequestShort, error) {
	return s.repo.GetReview(userID)
}
