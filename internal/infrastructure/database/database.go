package database

import (
	"fmt"
	"othnx/internal/domain/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(c entities.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DB.PostgresURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	} else {
		fmt.Println("Connected to database")
	}

	if err := db.AutoMigrate(&entities.User{}, &entities.RefreshToken{}, &entities.State{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entities.User{}, &entities.RefreshToken{}, &entities.State{})
}
