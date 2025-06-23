package domain

import "testing"

func TestComputeNameFingerprint(t *testing.T) {
	t.Run("should compute name fingerprint", func(t *testing.T) {
		name := "John Doe"
		fingerprint, err := GenerateNameFingerprint(name)
		if err != nil {
			t.Fatal(err)
		}
		if fingerprint != "320b8e6bef45211f0f57b618925f4193" {
			t.Fatal("expected fingerprint to be 320b8e6bef45211f0f57b618925f4193", fingerprint)
		}
	})
	t.Run("prevent naming collisions", func(t *testing.T) {
		name1 := "John Doe"
		fp1, err := GenerateNameFingerprint(name1)
		if err != nil {
			t.Fatal(err)
		}
		name2 := "John doé"
		fp2, err := GenerateNameFingerprint(name2)
		if err != nil {
			t.Fatal(err)
		}
		if fp1 != fp2 {
			t.Fatal("expected fingerprints to be the same", fp1, fp2)
		}
	})
	t.Run("2 different names should have different fingerprints", func(t *testing.T) {
		name1 := "John Doe"
		fp1, err := GenerateNameFingerprint(name1)
		if err != nil {
			t.Fatal(err)
		}
		name2 := "Jane Doe"
		fp2, err := GenerateNameFingerprint(name2)
		if err != nil {
			t.Fatal(err)
		}
		if fp1 == fp2 {
			t.Fatal("expected fingerprints to be different", fp1, fp2)
		}
	})
}

func TestUser(t *testing.T) {
	t.Run("should populate fields properly", func(t *testing.T) {
		user, err := NewUser("John Doe", "https://github.com/john-doe.png", "john.doe@example.com", "github")
		if err != nil {
			t.Fatal(err)
		}
		if user.ID == "" {
			t.Fatal("expected ID to be set")
		}
		if user.CreatedAt.IsZero() {
			t.Fatal("expected CreatedAt to be set")
		}
		if user.DeletedAt != nil {
			t.Fatal("expected DeletedAt to be nil")
		}
		if user.BlockedAt != nil {
			t.Fatal("expected BlockedAt to be nil")
		}
		if user.EarlyAdopter {
			t.Fatal("expected EarlyAdopter to be false")
		}
		if user.Name != "John Doe" {
			t.Fatal("expected name to be John Doe", user.Name)
		}
		if user.NameFingerprint != "320b8e6bef45211f0f57b618925f4193" {
			t.Fatal("expected name fingerprint to be 320b8e6bef45211f0f57b618925f4193", user.NameFingerprint)
		}
		if user.Email != "john.doe@example.com" {
			t.Fatal("expected email to be john.doe@example.com", user.Email)
		}
		if user.AvatarURL != "https://github.com/john-doe.png" {
			t.Fatal("expected avatar URL to be https://github.com/john-doe.png", user.AvatarURL)
		}
		if user.AuthMethod != "github" {
			t.Fatal("expected auth method to be github", user.AuthMethod)
		}
	})
	t.Run("should generate different users with different ids", func(t *testing.T) {
		user1, _ := NewUser("John Doe", "https://github.com/john-doe.png", "john.doe@example.com", "github")
		user2, _ := NewUser("John Doe", "https://github.com/john-doe.png", "john.doe@example.com", "github")
		if user1.ID == user2.ID {
			t.Fatal("expected users to have different ids", user1.ID, user2.ID)
		}
	})
}
