package service

import (
	"reviewers/internal/models"
	"reviewers/internal/repository"
)

type TeamService struct {
	repo *repository.TeamRepository
}

func NewTeamService(repo *repository.TeamRepository) *TeamService {
	return &TeamService{repo}
}

func (s *TeamService) GetTeam(name string) (*models.Team, error) {
	return s.repo.GetTeam(name)
}

func (s *TeamService) CreateTeam(newTeam *models.Team) error {
	return s.repo.CreateTeam(newTeam)
}

func (s *TeamService) GetReviewerIdsFromUserTeam(userID string) ([]*models.User, error) {
	return s.repo.GetReviewerIdsFromUserTeam(userID)
}
