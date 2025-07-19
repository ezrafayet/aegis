package integration_test_cases

import (
	"aegis/integration/integration_testkit"
	"aegis/internal/domain/entities"
	"aegis/pkg/jwtgen"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MiddlewareInternalAPI_NoKeyReturns401(t *testing.T) {
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

	req, err := http.NewRequest("POST", suite.Server.URL+"/auth/authorize-access-token", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorize", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 401, resp.StatusCode)
}

func MiddlewareInternalAPI_InvalidKeyReturns401(t *testing.T) {
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

	req, err := http.NewRequest("POST", suite.Server.URL+"/auth/authorize-access-token", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorize", "Bearer wrong_key")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 401, resp.StatusCode)
}
