package integration_test_cases

import (
	"aegis/integration/integration_testkit"
	"aegis/internal/domain/entities"
	"aegis/pkg/cookies"
	"aegis/pkg/jwtgen"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Login_NoATOrRT_Returns200AndShowsLoginPage(t *testing.T) {
	config := integration_testkit.GetBaseConfig()

	suite := integration_testkit.SetupTestSuite(t, config)
	defer suite.Teardown()

	resp, err := http.Get(suite.Server.URL + config.LoginPage.FullPath)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "<!DOCTYPE html>")
}

func Login_ValidAT_RedirectsToSuccessPage(t *testing.T) {
	config := integration_testkit.GetBaseConfig()

	suite := integration_testkit.SetupTestSuite(t, config)
	defer suite.Teardown()

	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.Metadata)
	require.NoError(t, err)
	accessToken, atExp, err := jwtgen.Generate(cClaims.ToMap(), time.Now(), 15, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", suite.Server.URL+config.LoginPage.FullPath, nil)
	require.NoError(t, err)

	accessCookie := cookies.NewAccessCookie(accessToken, atExp, suite.Config)
	req.AddCookie(&accessCookie)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	location := resp.Header.Get("Location")
	assert.Equal(t, config.App.RedirectAfterSuccess, location)
}

func Login_DisabledReturns404(t *testing.T) {
	config := integration_testkit.GetBaseConfig()
	config.LoginPage.Enabled = false

	suite := integration_testkit.SetupTestSuite(t, config)
	defer suite.Teardown()

	resp, err := http.Get(suite.Server.URL + config.LoginPage.FullPath)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
