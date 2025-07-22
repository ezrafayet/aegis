package config

import (
	"os"
	"testing"
)

func TestReadWithEnvReplacement(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_DB_URL", "postgres://user:pass@localhost:5432/auth")
	os.Setenv("TEST_JWT_SECRET", "my-secret-key")
	os.Setenv("TEST_APP_URL", "https://test.example.com")
	os.Setenv("TEST_GITHUB_CLIENT_ID", "github-client-123")
	os.Setenv("TEST_CORS_ORIGIN", "https://test.example.com")
	os.Setenv("PROTOCOL", "https")
	os.Setenv("HOST", "api.example.com")
	os.Setenv("PORT", "8080")
	defer func() {
		os.Unsetenv("TEST_DB_URL")
		os.Unsetenv("TEST_JWT_SECRET")
		os.Unsetenv("TEST_APP_URL")
		os.Unsetenv("TEST_GITHUB_CLIENT_ID")
		os.Unsetenv("TEST_CORS_ORIGIN")
		os.Unsetenv("PROTOCOL")
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
	}()

	// Create a temporary config file with various env: references
	configContent := `{
		"app": {
			"name": "Test App",
			"url": "${env:TEST_APP_URL}",
			"port": 8080,
			"cors_allowed_origins": ["${env:TEST_CORS_ORIGIN}", "http://localhost:3000"],
			"redirect_after_success": "mixed_${env:TEST_JWT_SECRET}_text",
			"redirect_after_error": "${env:PROTOCOL}://${env:HOST}:${env:PORT}"
		},
		"db": {
			"postgres_url": "${env:TEST_DB_URL}"
		},
		"jwt": {
			"secret": "${env:TEST_JWT_SECRET}",
			"access_token_expiration_minutes": 30,
			"refresh_token_expiration_days": 7
		},
		"auth": {
			"providers": {
				"github": {
					"enabled": true,
					"client_id": "${env:TEST_GITHUB_CLIENT_ID}",
					"client_secret": "static-github-secret"
				}
			}
		}
	}`

	tmpFile, err := os.CreateTemp("", "test-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Read the config
	config, err := Read(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Verify all replacements worked correctly
	if config.App.URL != "https://test.example.com" {
		t.Errorf("Expected app URL to be 'https://test.example.com', got '%s'", config.App.URL)
	}
	if config.DB.PostgresURL != "postgres://user:pass@localhost:5432/auth" {
		t.Errorf("Expected DB URL to be 'postgres://user:pass@localhost:5432/auth', got '%s'", config.DB.PostgresURL)
	}
	if config.JWT.Secret != "my-secret-key" {
		t.Errorf("Expected JWT secret to be 'my-secret-key', got '%s'", config.JWT.Secret)
	}
	if config.App.RedirectAfterSuccess != "mixed_my-secret-key_text" {
		t.Errorf("Expected mixed text to be interpolated, got '%s'", config.App.RedirectAfterSuccess)
	}
	if config.App.RedirectAfterError != "https://api.example.com:8080" {
		t.Errorf("Expected multiple variables to be replaced, got '%s'", config.App.RedirectAfterError)
	}
	if config.Auth.Providers.GitHub.ClientID != "github-client-123" {
		t.Errorf("Expected GitHub client ID to be 'github-client-123', got '%s'", config.Auth.Providers.GitHub.ClientID)
	}

	// Check slice handling
	if len(config.App.CorsAllowedOrigins) != 2 {
		t.Errorf("Expected 2 CORS origins, got %d", len(config.App.CorsAllowedOrigins))
	}
	if config.App.CorsAllowedOrigins[0] != "https://test.example.com" {
		t.Errorf("Expected first CORS origin to be 'https://test.example.com', got '%s'", config.App.CorsAllowedOrigins[0])
	}
	if config.App.CorsAllowedOrigins[1] != "http://localhost:3000" {
		t.Errorf("Expected second CORS origin to be 'http://localhost:3000', got '%s'", config.App.CorsAllowedOrigins[1])
	}
}
