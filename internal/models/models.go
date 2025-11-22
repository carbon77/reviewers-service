package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username string
	IsActive bool

	Teams        []Team        `gorm:"many2many:user_teams"`
	PullRequests []PullRequest `gorm:"foreignKey:AuthorID"`
	Reviews      []PullRequest `gorm:"many2many:pull_request_reviewers"`
}

type Team struct {
	TeamID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name   string

	Users []User `gorm:"many2many:user_teams"`
}

type UserTeam struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

type PullRequest struct {
	PullRequestID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	PullRequestName string
	Status          string
	CreatedAt       time.Time
	MergedAt        *time.Time

	AuthorID uuid.UUID
	Author   User `gorm:"foreignKey:AuthorID"`

	Reviewers []User `gorm:"many2many:pull_request_reviewers"`
}

type PullRequestReviewer struct {
	PullRequestID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;primaryKey"`
}
