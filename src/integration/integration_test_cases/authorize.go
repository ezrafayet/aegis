package integration_test_cases

import (
	"aegis/internal/domain/entities"
	"aegis/pkg/apperrors"
	"aegis/pkg/jwtgen"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"aegis/integration/integration_testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AuthorizeRequest struct {
	AccessToken string   `json:"access_token"`
	Roles       []string `json:"authorized_roles"`
}

func Authorize_EmptyTokenReturns401(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	reqBody := AuthorizeRequest{
		AccessToken: "",
		Roles:       []string{"user"},
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.Server.URL+"/auth/authorize-access-token", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 401, resp.StatusCode)
	body := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, apperrors.ErrAccessTokenInvalid.Error(), body["error"])
	assert.Equal(t, false, body["authorized"])
}

func Authorize_ExpiredTokenReturns401(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user and generate an expired token
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.MetadataPublic)
	require.NoError(t, err)
	// Generate token with past expiration
	expiredToken, _, err := jwtgen.Generate(cClaims.ToMap(), time.Now().Add(-time.Hour*24), 15, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	reqBody := AuthorizeRequest{
		AccessToken: expiredToken,
		Roles:       []string{"user"},
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.Server.URL+"/auth/authorize-access-token", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	body := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, apperrors.ErrAccessTokenExpired.Error(), body["error"])
	assert.Equal(t, false, body["authorized"])
}

func Authorize_MalformedTokenReturns401(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	reqBody := AuthorizeRequest{
		AccessToken: "not.a.valid.jwt.token",
		Roles:       []string{"user"},
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.Server.URL+"/auth/authorize-access-token", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	body := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, apperrors.ErrAccessTokenInvalid.Error(), body["error"])
	assert.Equal(t, false, body["authorized"])
}

func Authorize_NoRolesReturns401(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user and generate a valid token
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.MetadataPublic)
	require.NoError(t, err)
	validToken, _, err := jwtgen.Generate(cClaims.ToMap(), time.Now(), 15, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	reqBody := AuthorizeRequest{
		AccessToken: validToken,
		Roles:       []string{}, // Empty roles array
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.Server.URL+"/auth/authorize-access-token", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	body := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, apperrors.ErrNoRoles.Error(), body["error"])
	assert.Equal(t, false, body["authorized"])
}

func Authorize_UserDoesNotHaveRequiredRoleReturns401(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user with only "user" role
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.MetadataPublic)
	require.NoError(t, err)
	validToken, _, err := jwtgen.Generate(cClaims.ToMap(), time.Now(), 15, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	reqBody := AuthorizeRequest{
		AccessToken: validToken,
		Roles:       []string{"admin"}, // User doesn't have admin role
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.Server.URL+"/auth/authorize-access-token", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	body := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, apperrors.ErrUnauthorizedRole.Error(), body["error"])
	assert.Equal(t, false, body["authorized"])
}

func Authorize_UserHasRequiredRoleReturns200(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user with "user" role
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.MetadataPublic)
	require.NoError(t, err)
	validToken, _, err := jwtgen.Generate(cClaims.ToMap(), time.Now(), 15, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	reqBody := AuthorizeRequest{
		AccessToken: validToken,
		Roles:       []string{"user"}, // User has this role
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.Server.URL+"/auth/authorize-access-token", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var response map[string]bool
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)
	assert.True(t, response["authorized"])
}

func Authorize_UserHasAnyRoleReturns200(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user with "user" role
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.MetadataPublic)
	require.NoError(t, err)
	validToken, _, err := jwtgen.Generate(cClaims.ToMap(), time.Now(), 15, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	reqBody := AuthorizeRequest{
		AccessToken: validToken,
		Roles:       []string{"payments", "any"}, // "any" should authorize any user
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.Server.URL+"/auth/authorize-access-token", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var response map[string]bool
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)
	assert.True(t, response["authorized"])
}
