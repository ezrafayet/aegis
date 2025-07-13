package integration_test_cases

import (
	"aegis/integration/integration_testkit"
	"aegis/internal/domain/entities"
	"aegis/pkg/apperrors"
	"aegis/pkg/cookies"
	"aegis/pkg/jwtgen"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func HardRefresh_MustRefresh_ValidTokens(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	// Create valid access token
	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.Metadata)
	require.NoError(t, err)
	accessToken, atExp, err := jwtgen.Generate(cClaims.ToMap(), time.Now(), 15, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	// Create valid refresh token
	refreshTokenEntity, _, err := entities.NewRefreshToken(user, "device-fingerprint", suite.Config)
	require.NoError(t, err)
	refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)

	// Create request with both tokens
	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/refresh", nil)
	require.NoError(t, err)

	accessCookie := cookies.NewAccessCookie(accessToken, atExp, suite.Config)
	refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
	req.AddCookie(&accessCookie)
	req.AddCookie(&refreshCookie)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify new cookies are set
	cookies := resp.Cookies()
	var newAccessToken, newRefreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			newAccessToken = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			newRefreshToken = cookie.Value
		}
	}
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	// assert.NotEqual(t, accessToken, newAccessToken) // todo: fix it, should be a new AT
	assert.NotEqual(t, refreshTokenEntity.Token, newRefreshToken)

	// Verify old refresh token is deleted
	var count int64
	err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func HardRefresh_MustRefresh_EmptyAT(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	// Create valid refresh token
	refreshTokenEntity, _, err := entities.NewRefreshToken(user, "device-fingerprint", suite.Config)
	require.NoError(t, err)
	refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)

	// Create request with empty access token and valid refresh token
	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/refresh", nil)
	require.NoError(t, err)

	refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
	req.AddCookie(&refreshCookie)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify new cookies are set
	cookies := resp.Cookies()
	var newAccessToken, newRefreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			newAccessToken = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			newRefreshToken = cookie.Value
		}
	}
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, refreshTokenEntity.Token, newRefreshToken)

	// Verify old refresh token is deleted
	var count int64
	err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func HardRefresh_MustRefresh_ExpiredAT(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	// Create expired access token (1 minute expiration, created 2 minutes ago)
	cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.Metadata)
	require.NoError(t, err)
	accessToken, atExp, err := jwtgen.Generate(cClaims.ToMap(), time.Now().Add(-2*time.Minute), 1, "TestApp", suite.Config.JWT.Secret)
	require.NoError(t, err)

	// Create valid refresh token
	refreshTokenEntity, _, err := entities.NewRefreshToken(user, "device-fingerprint", suite.Config)
	require.NoError(t, err)
	refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)

	// Create request with expired access token and valid refresh token
	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/refresh", nil)
	require.NoError(t, err)

	accessCookie := cookies.NewAccessCookie(accessToken, atExp, suite.Config)
	refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
	req.AddCookie(&accessCookie)
	req.AddCookie(&refreshCookie)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify new cookies are set
	cookies := resp.Cookies()
	var newAccessToken, newRefreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			newAccessToken = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			newRefreshToken = cookie.Value
		}
	}
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, accessToken, newAccessToken)
	assert.NotEqual(t, refreshTokenEntity.Token, newRefreshToken)

	// Verify old refresh token is deleted
	var count int64
	err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func HardRefresh_MustRefresh_MalformedAT(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	// Create malformed access token
	malformedAccessToken := "invalid.jwt.token"

	// Create valid refresh token
	refreshTokenEntity, _, err := entities.NewRefreshToken(user, "device-fingerprint", suite.Config)
	require.NoError(t, err)
	refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)

	// Create request with malformed access token and valid refresh token
	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/refresh", nil)
	require.NoError(t, err)

	accessCookie := cookies.NewAccessCookie(malformedAccessToken, time.Now().Add(time.Hour).Unix(), suite.Config)
	refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
	req.AddCookie(&accessCookie)
	req.AddCookie(&refreshCookie)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify new cookies are set
	cookies := resp.Cookies()
	var newAccessToken, newRefreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			newAccessToken = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			newRefreshToken = cookie.Value
		}
	}
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)

	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, malformedAccessToken, newAccessToken)
	assert.NotEqual(t, refreshTokenEntity.Token, newRefreshToken)

	// Verify old refresh token is deleted
	var count int64
	err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func HardRefresh_MustNotRefresh_EmptyRT(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create request without refresh token
	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/refresh", nil)
	require.NoError(t, err)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Verify error response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	require.NoError(t, err)
	// todo: could consider another error
	assert.Contains(t, errorResponse["error"], apperrors.ErrGeneric.Error())
}

func HardRefresh_MustNotRefresh_ExpiredRT(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create a user
	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
	require.NoError(t, err)
	user = suite.CreateUser(t, user, []string{"user"})

	// Create expired refresh token
	refreshTokenEntity := entities.RefreshToken{
		UserID:            user.ID,
		CreatedAt:         time.Now().Add(-8 * 24 * time.Hour), // 8 days ago
		ExpiresAt:         time.Now().Add(-1 * 24 * time.Hour), // 1 day ago
		Token:             "expired_refresh_token",
		DeviceFingerprint: "device-fingerprint",
	}
	refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)

	// Create request with expired refresh token
	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/refresh", nil)
	require.NoError(t, err)

	refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
	req.AddCookie(&refreshCookie)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Verify error response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	require.NoError(t, err)
	assert.Contains(t, errorResponse["error"], apperrors.ErrRefreshTokenExpired.Error())
}

func HardRefresh_MustNotRefresh_MalformedRT(t *testing.T) {
	suite := integration_testkit.SetupTestSuite(t, integration_testkit.GetBaseConfig())
	defer suite.Teardown()

	// Create request with malformed refresh token
	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/refresh", nil)
	require.NoError(t, err)

	refreshCookie := cookies.NewRefreshCookie("malformed_refresh_token", time.Now().Add(time.Hour).Unix(), suite.Config)
	req.AddCookie(&refreshCookie)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Verify error response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	require.NoError(t, err)
	// todo: could consider another error
	assert.Contains(t, errorResponse["error"], apperrors.ErrGeneric.Error())
}
