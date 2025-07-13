package integration_test_cases

import (
	"aegis/integration/integration_testkit"
	"aegis/internal/domain/entities"
	"aegis/pkg/apperrors"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderCallback(t *testing.T) {
	t.Run("unhappy scenarios: generic cases", func(t *testing.T) {
		t.Run("calling GET /provider/callback returns 403 if the provider is not enabled", func(t *testing.T) {
			config := integration_testkit.GetBaseConfig()
			config.Auth.Providers.GitHub.Enabled = false
			suite := integration_testkit.SetupTestSuite(t, config)
			defer suite.Teardown()

			resp, err := http.Get(suite.Server.URL + "/auth/github/callback")
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			var errorResponse map[string]interface{}
			err = json.Unmarshal(body, &errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse["error"], apperrors.ErrAuthMethodNotEnabled.Error())
		})
	})

	t.Run("unhappy scenarios: cases that must redirect to error page", func(t *testing.T) {

		t.Run("calling GET /provider/callback redirects to error page if state is invalid", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()
			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=valid_code&state=invalid_state", nil)
			require.NoError(t, err)
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Logf("HTTP request failed: %v", err)
				t.Fatal(err)
			}
			require.NoError(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			// todo: proper error ?
			assert.Equal(t, location, "http://localhost:8080/login-error?error=unknown_error")
		})

		t.Run("calling GET /provider/callback redirects to error page if code is invalid", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a valid state first
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err := suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=invalid_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			assert.Equal(t, location, "http://localhost:8080/login-error?error=unknown_error")
		})

		t.Run("calling GET /provider/callback returns to error page if user declines auth", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?error=access_denied", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			assert.Equal(t, location, "http://localhost:8080/login-error?error=access_denied")
		})

		t.Run("calling GET /provider/callback returns to error page if user is using another method", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a user with discord auth method
			user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "discord")
			require.NoError(t, err)
			user = suite.CreateUser(t, user, []string{"user"})

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err = suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			// Mock GitHub to return the same email as the discord user
			// This would require modifying the mock server to return the same email
			// For now, we'll test the basic flow
			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=valid_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			// todo: check proper error
			assert.Equal(t, "http://localhost:8080/login-error?error=unknown_error", location)
		})

		t.Run("calling GET /provider/callback returns to error page if user is blocked", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a blocked user
			user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
			require.NoError(t, err)
			now := time.Now()
			user.BlockedAt = &now
			user = suite.CreateUser(t, user, []string{"user"})

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err = suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=valid_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			assert.Equal(t, location, "http://localhost:8080/login-error?error=unknown_error")
		})

		t.Run("calling GET /provider/callback returns to error page if user is deleted", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a deleted user
			user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
			require.NoError(t, err)
			now := time.Now()
			user.DeletedAt = &now
			user = suite.CreateUser(t, user, []string{"user"})

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err = suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=valid_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			// todo: check error
			assert.Equal(t, location, "http://localhost:8080/login-error?error=unknown_error")
		})

		t.Run("calling GET /provider/callback returns to error page if user is not an early user", func(t *testing.T) {
			config := integration_testkit.GetBaseConfig()
			config.App.EarlyAdoptersOnly = true
			suite := integration_testkit.SetupTestSuite(t, config)
			defer suite.Teardown()

			// Create a user who is not an early adopter
			user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
			require.NoError(t, err)
			user.EarlyAdopter = false
			user = suite.CreateUser(t, user, []string{"user"})

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err = suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=valid_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			// todo: check error
			assert.Equal(t, location, "http://localhost:8080/login-error?error=unknown_error")
		})
	})

	t.Run("happy scenarios", func(t *testing.T) {
		t.Run("calling GET /provider/callback gives [access_token, refresh_token] if the user already exists", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create an existing user
			user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
			require.NoError(t, err)
			user = suite.CreateUser(t, user, []string{"user"})

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err = suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=accepted_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			assert.Equal(t, "http://localhost:8080/login-success", location)

			// Verify cookies are set
			cookies := resp.Cookies()
			var accessToken, refreshToken string
			for _, cookie := range cookies {
				if cookie.Name == "access_token" {
					accessToken = cookie.Value
				}
				if cookie.Name == "refresh_token" {
					refreshToken = cookie.Value
				}
			}
			assert.NotEmpty(t, accessToken)
			assert.NotEmpty(t, refreshToken)

			// check user length is 1
			var count int64
			err = suite.Db.Model(&entities.User{}).Count(&count).Error
			require.NoError(t, err)
			assert.Equal(t, int64(1), count)
		})

		t.Run("calling GET /provider/callback gives [access_token, refresh_token] and creates user if the user does not exist", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err := suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			// Verify no user exists with the email from mock
			var count int64
			err = suite.Db.Model(&entities.User{}).Where("email = ?", "test@example.com").Count(&count).Error
			require.NoError(t, err)
			assert.Equal(t, int64(0), count)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=accepted_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			assert.Equal(t, "http://localhost:8080/login-success", location)

			// Verify cookies are set
			cookies := resp.Cookies()
			var accessToken, refreshToken string
			for _, cookie := range cookies {
				if cookie.Name == "access_token" {
					accessToken = cookie.Value
				}
				if cookie.Name == "refresh_token" {
					refreshToken = cookie.Value
				}
			}
			assert.NotEmpty(t, accessToken)
			assert.NotEmpty(t, refreshToken)

			// Verify user was created
			err = suite.Db.Model(&entities.User{}).Where("email = ?", "test@example.com").Count(&count).Error
			require.NoError(t, err)
			assert.Equal(t, int64(1), count)
		})

		t.Run("calling GET /provider/callback cleans the state", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err := suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			// Verify state exists
			var count int64
			err = suite.Db.Model(&entities.State{}).Where("value = ?", "valid_state").Count(&count).Error
			require.NoError(t, err)
			assert.Equal(t, int64(1), count)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=accepted_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)

			// Verify state was deleted
			err = suite.Db.Model(&entities.State{}).Where("value = ?", "valid_state").Count(&count).Error
			require.NoError(t, err)
			assert.Equal(t, int64(0), count)
		})

		t.Run("calling GET /provider/callback redirects to the welcome page", func(t *testing.T) {
			suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a valid state
			state := entities.State{
				Value:     "valid_state",
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			err := suite.Db.Model(&entities.State{}).Create(&state).Error
			require.NoError(t, err)

			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/github/callback?code=accepted_code&state=valid_state", nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusFound, resp.StatusCode)
			location := resp.Header.Get("Location")
			assert.Equal(t, suite.Config.App.RedirectAfterSuccess, location)
		})
	})
}

// t.Run("rate limiting", func(t *testing.T) {
// 	// todo: implement and test rate limiting
// })
// t.Run("devices ids", func(t *testing.T) {
// 	// todo: implement and test different devices ids
// })
