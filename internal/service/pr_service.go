package service

import (
	"log/slog"
	"math/rand/v2"
	"reviewers/internal/errs"
	"reviewers/internal/models"
	"reviewers/internal/repository"
	"slices"
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

	reviewers, err := s.teamService.GetReviewerIdsFromUserTeam(pr.AuthorID)
	if err != nil {
		return err
	}
	for _, r := range reviewers {
		slog.Info("reviewer: %+v", r)
	}

	pr.Reviewers = getRandomReviewers(pr.ID, reviewers, 2)
	for _, reviewer := range pr.Reviewers {
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer.UserID)
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
	return s.repo.Get(pullRequestID)
}

func (s *PRService) Reassign(pullRequestID, oldReviewerID string) (*models.PullRequest, error) {
	pr, err := s.repo.Get(pullRequestID)
	if err != nil || pr == nil {
		return nil, err
	}

	// Check if PR is already merged
	if pr.Status == models.StatusMerged {
		return nil, errs.PullRequestMerged
	}

	// Check if old reviewer is not assigned
	oldReviewerIdx := slices.Index(pr.AssignedReviewers, oldReviewerID)
	if oldReviewerIdx == -1 {
		return nil, errs.NotAssigned
	}

	reviewers, err := s.teamService.GetReviewerIdsFromUserTeam(pr.AuthorID, pr.AssignedReviewers...)
	if err != nil {
		return nil, err
	}

	if len(reviewers) == 0 {
		return nil, errs.NoCandidate
	}

	newReviewer := getRandomReviewers(pullRequestID, reviewers, 1)[0]
	pr.Reviewers[oldReviewerIdx] = newReviewer

	pr.AssignedReviewers = make([]string, 0, len(pr.Reviewers))
	for _, reviewer := range pr.Reviewers {
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer.UserID)
	}

	err = s.repo.Save(pr)
	return pr, err
}

func getRandomReviewers(pullRequestID string, reviewers []*models.User, count int) []models.PullRequestReviewer {
	shuffled := make([]models.PullRequestReviewer, 0, len(reviewers))
	for _, user := range reviewers {
		shuffled = append(shuffled, models.PullRequestReviewer{
			UserID:        user.ID,
			PullRequestID: pullRequestID,
		})
	}

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:min(len(shuffled), count)]
}
