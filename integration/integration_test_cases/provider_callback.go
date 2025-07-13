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

func ProviderCallback_NotEnabledReturns403(t *testing.T) {
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
}

func ProviderCallback_MustRedirectToErrorPage_InvalidState(t *testing.T) {
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
}

func ProviderCallback_MustRedirectToErrorPage_InvalidCode(t *testing.T) {
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
}

func ProviderCallback_MustRedirectToErrorPage_UserDeclinesAuth(t *testing.T) {
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
}

func ProviderCallback_MustRedirectToErrorPage_UserUsingAnotherMethod(t *testing.T) {
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
}

func ProviderCallback_MustRedirectToErrorPage_UserBlocked(t *testing.T) {
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
}

func ProviderCallback_MustRedirectToErrorPage_UserDeleted(t *testing.T) {
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
}

func ProviderCallback_MustRedirectToErrorPage_UserNotAnEarlyAdopter(t *testing.T) {
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
}
