package service

import (
	"math/rand/v2"
	"reviewers/internal/models"
	"reviewers/internal/repository"
	"time"
)

type PRService struct {
	repo        *repository.PRRepository
	teamService *TeamService
}

func NewPRService(repo *repository.PRRepository, teamService *TeamService) *PRService {
	return &PRService{repo, teamService}
}

func (s *PRService) Create(pr *models.PullRequest) error {
	pr.CreatedAt = time.Now()
	pr.Status = models.StatusOpen

	ids, err := s.teamService.GetReviewerIdsFromUserTeam(pr.AuthorID)
	if err != nil {
		return err
	}

	shuffled := make([]*models.User, len(ids))
	copy(shuffled, ids)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	assignedIds := shuffled[:min(len(shuffled), 2)]
	var assignedUsers []models.User

	for _, user := range assignedIds {
		assignedUsers = append(assignedUsers, *user)
	}

	pr.Reviewers = assignedUsers
	for _, reviewer := range pr.Reviewers {
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer.ID)
	}

	return s.repo.Create(pr)
}

func (s *PRService) Merge(pullRequestID string) (*models.PullRequest, error) {
	pr, err := s.repo.Get(pullRequestID)
	if err != nil || pr == nil {
		return nil, err
	}

	if pr.Status == models.StatusMerged {
		return pr, nil
	}

	err = s.repo.Merge(pullRequestID)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
