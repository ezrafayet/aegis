package repositories

import (
	"aegis/internal/domain/entities"
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
		db.AutoMigrate(&entities.User{}, &entities.Role{})
		userRepository := NewUserRepository(db)
		err = userRepository.CreateUser(entities.User{
			ID:              "123",
			CreatedAt:       time.Now(),
			Name:            "Test User",
			NameFingerprint: "test_fingerprint",
			Email:           "test@test.com",
			MetadataPublic:  "{}",
			AuthMethod:      "test",
		}, []entities.Role{entities.NewRole("123", "user")})
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
		if len(user.Roles) != 1 {
			t.Fatal("expected user to have 1 role", len(user.Roles))
		}
		if user.Roles[0].Value != "user" {
			t.Fatal("expected user to have role user", user.Roles[0].Value)
		}
	})
	t.Run("Forbid email collision", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.Role{})
		userRepository := NewUserRepository(db)
		err = userRepository.CreateUser(entities.User{
			ID:              "123",
			CreatedAt:       time.Now(),
			Name:            "Test User1",
			NameFingerprint: "test_fingerprint1",
			Email:           "test@test.com",
			MetadataPublic:  "{}",
			AuthMethod:      "test",
		}, []entities.Role{entities.NewRole("123", "user")})
		if err != nil {
			t.Fatal(err)
		}
		err = userRepository.CreateUser(entities.User{
			ID:              "456",
			CreatedAt:       time.Now(),
			Name:            "Test User2",
			NameFingerprint: "test_fingerprint2",
			Email:           "test@test.com",
			MetadataPublic:  "{}",
			AuthMethod:      "test",
		}, []entities.Role{entities.NewRole("456", "user")})
		if err == nil {
			t.Fatal(err)
		}
	})

	// Now authorized
	//
	//	t.Run("Forbid name collision", func(t *testing.T) {
	//		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//		db.AutoMigrate(&entities.User{}, &entities.Role{})
	//		userRepository := NewUserRepository(db)
	//		err = userRepository.CreateUser(entities.User{
	//			ID:              "123",
	//			CreatedAt:       time.Now(),
	//			Name:            "Test User1",
	//			NameFingerprint: "test_fingerprint1",
	//			Email:           "test1@test.com",
	//			MetadataPublic:  "{}",
	//			AuthMethod:      "test",
	//		}, []entities.Role{entities.NewRole("123", "user")})
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//		err = userRepository.CreateUser(entities.User{
	//			ID:              "456",
	//			CreatedAt:       time.Now(),
	//			Name:            "Test User1",
	//			NameFingerprint: "test_fingerprint1",
	//			Email:           "test2@test.com",
	//			MetadataPublic:  "{}",
	//			AuthMethod:      "test",
	//		}, []entities.Role{entities.NewRole("456", "user")})
	//		if err == nil {
	//			t.Fatal(err)
	//		}
	//	})
}

func TestUserRepository_Roles(t *testing.T) {
	t.Run("should return roles for a user", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.Role{})
		userRepository := NewUserRepository(db)
		err = userRepository.CreateUser(entities.User{
			ID:              "123",
			CreatedAt:       time.Now(),
			Name:            "Test User",
			NameFingerprint: "test_fingerprint",
			Email:           "test@test.com",
			MetadataPublic:  "{}",
			AuthMethod:      "test",
		}, []entities.Role{entities.NewRole("123", "user")})
		if err != nil {
			t.Fatal(err)
		}
		db.Create(entities.NewRole("123", "admin"))
		user, err := userRepository.GetUserByEmail("test@test.com")
		if err != nil {
			t.Fatal(err)
		}
		if len(user.Roles) != 2 {
			t.Fatal("expected user to have 2 roles", len(user.Roles))
		}
		if user.Roles[0].Value != "admin" {
			t.Fatal("expected user to have role user", user.Roles[0].Value)
		}
		if user.Roles[1].Value != "user" {
			t.Fatal("expected user to have role admin", user.Roles[1].Value)
		}
	})
	t.Run("should not break if user has no roles", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.Role{})
		userRepository := NewUserRepository(db)
		err = userRepository.CreateUser(entities.User{
			ID:              "123",
			CreatedAt:       time.Now(),
			Name:            "Test User",
			NameFingerprint: "test_fingerprint",
			Email:           "test@test.com",
			MetadataPublic:  "{}",
			AuthMethod:      "test",
		}, []entities.Role{})
		if err != nil {
			t.Fatal(err)
		}
		user, err := userRepository.GetUserByEmail("test@test.com")
		if err != nil {
			t.Fatal(err)
		}
		if len(user.Roles) != 0 {
			t.Fatal("expected user to have 0 roles", len(user.Roles))
		}
	})
}
