package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"user_id" gorm:"column:user_id;primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	IsActive bool   `json:"is_active"`
	TeamID   string `json:"-"`
}

type Team struct {
	ID      string `json:"-" gorm:"column:team_id;primaryKey"`
	Name    string `json:"team_name" gorm:"column:name;unique;not null"`
	Members []User `json:"members" gorm:"foreignKey:TeamID"`
}

type PullRequest struct {
	ID        string     `json:"pull_request_id" gorm:"column:pull_request_id;primaryKey"`
	Name      string     `json:"pull_request_name" gorm:"column:pull_request_name"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	MergedAt  *time.Time `json:"merged_at"`

	AuthorID string `json:"author_id"`
	Author   User   `json:"-" gorm:"foreignKey:AuthorID"`

	Reviewers         []PullRequestReviewer `json:"-" gorm:"foreignKey:PullRequestID"`
	AssignedReviewers []string              `json:"assigned_reviewers" gorm:"-"`
}

func (pr *PullRequest) AfterFind(tx *gorm.DB) (err error) {
	if len(pr.Reviewers) > 0 {
		pr.AssignedReviewers = make([]string, 0, len(pr.Reviewers))
		for _, reviewer := range pr.Reviewers {
			pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer.UserID)
		}
	}
	return nil
}

type PullRequestReviewer struct {
	PullRequestID string `json:"pull_request_id" gorm:"column:pull_request_id;primaryKey"`
	UserID        string `json:"user_id" gorm:"column:user_id;primaryKey"`
}

const StatusOpen = "OPEN"
const StatusMerged = "MERGED"

type PullRequestShort struct {
	ID       string `json:"pull_request_id" gorm:"column:pull_request_id"`
	Name     string `json:"pull_request_name" gorm:"column:pull_request_name"`
	AuthorID string `json:"author_id" gorm:"column:author_id"`
	Status   string `json:"status" gorm:"column:status"`
}
