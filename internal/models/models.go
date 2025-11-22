package models

import (
	"time"
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

	Reviewers []User `json:"-" gorm:"many2many:pull_request_reviewers"`
}

const StatusOpen = "OPEN"
const StatusMerged = "MERGED"
