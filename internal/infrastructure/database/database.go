package database

import (
	"fmt"
	"othnx/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(c domain.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DB.PostgresURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	} else {
		fmt.Println("Connected to database")
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}, &domain.State{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}, &domain.State{})
}
