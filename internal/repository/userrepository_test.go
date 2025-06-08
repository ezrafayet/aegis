package repository

import (
	"aegix/internal/domain"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	t.Run("Most basic test", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{})
		userRepository := NewUserRepository(db)
		err = userRepository.CreateUser(domain.User{
			ID:              "123",
			CreatedAt:       time.Now(),
			Name:            "Test User",
			NameFingerprint: "test_fingerprint",
			Email:           "test@test.com", 
			Metadata:        "{}",
			AuthMethod:      "test",
		})
		if err != nil {
			t.Fatal(err)
		}
		user, err := userRepository.GetUserByEmail("test@test.com")
		if err != nil {
			t.Fatal(err)
		}
		if user.ID != "123" {
			t.Fatal("expected user id to be 123", user.ID)
		}
		if user.Name != "Test User" {
			t.Fatal("expected user name to be Test User", user.Name)
		}
		if user.NameFingerprint != "test_fingerprint" {
			t.Fatal("expected user name fingerprint to be test_fingerprint", user.NameFingerprint)
		}
		if user.Email != "test@test.com" {
			t.Fatal("expected user email to be test@test.com", user.Email)
		}
		if user.AuthMethod != "test" {
			t.Fatal("expected user auth method to be test", user.AuthMethod)
		}
	})
	t.Run("Forbid email collision", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{})
		userRepository := NewUserRepository(db)
		err = userRepository.CreateUser(domain.User{
			ID:              "123",
			CreatedAt:       time.Now(),
			Name:            "Test User1",
			NameFingerprint: "test_fingerprint1",
			Email:           "test@test.com", 
			Metadata:        "{}",
			AuthMethod:      "test",
		})
		if err != nil {
			t.Fatal(err)
		}
		err = userRepository.CreateUser(domain.User{
			ID:              "456",
			CreatedAt:       time.Now(),
			Name:            "Test User2",
			NameFingerprint: "test_fingerprint2",
			Email:           "test@test.com", 
			Metadata:        "{}",
			AuthMethod:      "test",
		})
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("Forbid name collision", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{})
		userRepository := NewUserRepository(db)
		err = userRepository.CreateUser(domain.User{
			ID:              "123",
			CreatedAt:       time.Now(),
			Name:            "Test User1",
			NameFingerprint: "test_fingerprint1",
			Email:           "test1@test.com", 
			Metadata:        "{}",
			AuthMethod:      "test",
		})
		if err != nil {
			t.Fatal(err)
		}
		err = userRepository.CreateUser(domain.User{
			ID:              "456",
			CreatedAt:       time.Now(),
			Name:            "Test User1",
			NameFingerprint: "test_fingerprint1",
			Email:           "test2@test.com", 
			Metadata:        "{}",
			AuthMethod:      "test",
		})
		if err == nil {
			t.Fatal(err)
		}
	})
}
