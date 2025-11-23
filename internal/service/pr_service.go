package service

import (
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

	pr.Reviewers = getRandomReviewers(reviewers, 2)
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
	for _, user := range pr.Reviewers {
		pr.AssignedReviewers = append(pr.AssignedReviewers, user.ID)
	}

	if !slices.Contains(pr.AssignedReviewers, oldReviewerID) {
		return nil, errs.NotAssigned
	}

	reviewers, err := s.teamService.GetReviewerIdsFromUserTeam(pr.AuthorID, pr.AssignedReviewers...)
	if err != nil {
		return nil, err
	}

	if len(reviewers) == 0 {
		return nil, errs.NoCandidate
	}

	newReviewer := getRandomReviewers(reviewers, 1)[0]
	oldReviewerIdx := slices.Index(pr.AssignedReviewers, oldReviewerID)

	pr.Reviewers[oldReviewerIdx] = newReviewer
	pr.AssignedReviewers = make([]string, 0, len(pr.Reviewers))
	for _, user := range pr.Reviewers {
		pr.AssignedReviewers = append(pr.AssignedReviewers, user.ID)
	}
	err = s.repo.Save(pr)
	return pr, err
}

func getRandomReviewers(reviewers []*models.User, count int) []models.User {
	shuffled := make([]*models.User, len(reviewers))
	copy(shuffled, reviewers)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	assignedReviewers := shuffled[:min(len(shuffled), count)]
	var assignedUsers []models.User

	for _, user := range assignedReviewers {
		assignedUsers = append(assignedUsers, *user)
	}
	return assignedUsers
}
