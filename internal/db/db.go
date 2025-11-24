package db

import (
	"fmt"
	"reviewers/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
}
