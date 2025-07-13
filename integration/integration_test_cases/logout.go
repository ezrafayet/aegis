package integration_test_cases

import (
	"aegis/internal/domain/entities"
	"aegis/pkg/cookies"
	"net/http"
	"testing"
	"time"

	"aegis/integration/integration_testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Logout_SetsZeroCookies(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()
	resp, err := http.Get(suite.Server.URL + "/auth/logout")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cookies := resp.Cookies()
	assert.Len(t, cookies, 2)
	require.NotNil(t, cookies[0])
	assert.Equal(t, "", cookies[0].Value)
	assert.True(t, cookies[0].Expires.Before(time.Now()) || cookies[0].Expires.Equal(time.Unix(0, 0)))
	require.NotNil(t, cookies[1])
	assert.Equal(t, "", cookies[1].Value)
	assert.True(t, cookies[1].Expires.Before(time.Now()) || cookies[1].Expires.Equal(time.Unix(0, 0)))
}

func Logout_WithoutRefreshTokenDoesNotBreak(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()
	resp, err := http.Get(suite.Server.URL + "/auth/logout")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Logout_DeletesRefreshToken(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
		defer suite.Teardown()

		// Create a user
		user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
		require.NoError(t, err)
		user = suite.CreateUser(t, user, []string{"user"})

		refreshTokenEntity, _, err := entities.NewRefreshToken(user, "refresh_token", suite.Config)
		require.NoError(t, err)
		refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)

		// Verify the refresh token exists in the database
		var count int64
		err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Create a request with the refresh token cookie
		req, err := http.NewRequest("GET", suite.Server.URL+"/auth/logout", nil)
		require.NoError(t, err)

		refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
		req.AddCookie(&refreshCookie)

		// Make the request
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the refresh token was deleted from the database
		err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
}
