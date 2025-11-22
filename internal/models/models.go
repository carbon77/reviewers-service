package models

import (
	"time"
)

type User struct {
	ID       string `json:"user_id" gorm:"column:user_id;primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	IsActive bool   `json:"is_active"`

	Teams        []Team        `json:"-" gorm:"many2many:user_teams"`
	PullRequests []PullRequest `json:"-" gorm:"foreignKey:AuthorID"`
	Reviews      []PullRequest `json:"-" gorm:"many2many:pull_request_reviewers"`
}

type Team struct {
	ID   string `json:"team_id" gorm:"column:team_id;primaryKey"`
	Name string `json:"team_name" gorm:"unique;not null"`

	Members []User `json:"members" gorm:"many2many:user_teams"`
}

type PullRequest struct {
	PullRequestID   string `gorm:"primaryKey"`
	PullRequestName string
	Status          string
	CreatedAt       time.Time
	MergedAt        *time.Time

	AuthorID string
	Author   User `gorm:"foreignKey:AuthorID"`

	Reviewers []User `gorm:"many2many:pull_request_reviewers"`
}

type PullRequestReviewer struct {
	PullRequestID string `gorm:"primaryKey"`
	UserID        string `gorm:"primaryKey"`
}
