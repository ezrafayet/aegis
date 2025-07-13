package integration_test_cases

import (
	"aegis/integration/integration_testkit"
	"aegis/internal/domain/entities"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderCallback(t *testing.T) {
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
