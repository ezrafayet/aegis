package integration

import (
	"aegis/integration-tests/testkit"
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

func AreAccessTokenClaimsValid(accessToken string, expectedClaims map[string]interface{}) []error {
	// todo
	return nil
}

func TestRefresh_HardRefresh(t *testing.T) {
	t.Run("unhappy scenarios: user should not be refreshed", func(t *testing.T) {
		t.Run("user does not exist", func(t *testing.T) {})
		t.Run("user is deleted", func(t *testing.T) {})
		t.Run("user is blocked", func(t *testing.T) {})
		t.Run("user is not an early user", func(t *testing.T) {})
	})

	t.Run("calling GET /refresh with [valid access_token, valid refresh_token] returns a new [access_token, refresh_token]", func(t *testing.T) {
		suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
	})

	t.Run("calling GET /refresh with [invalid access_token, valid refresh_token] returns a new [access_token, refresh_token]", func(t *testing.T) {
		t.Run("access_token is expired", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
		})

		t.Run("access_token is malformed", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
		})

		t.Run("access_token is empty", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
		})
	})

	t.Run("calling GET /refresh with [invalid access_token, invalid refresh_token] returns 401", func(t *testing.T) {
		t.Run("refresh_token is expired", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
		})

		t.Run("refresh_token is malformed", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
		})

		t.Run("refresh_token is empty", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
		})
	})
}

func TestRefresh_SoftRefresh(t *testing.T) {
	t.Run("calling a route with check & soft refresh, with [valid access_token, valid refresh_token] does not update the tokens", func(t *testing.T) {
		suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
		req, err := http.NewRequest("GET", suite.Server.URL+"/auth/me", nil)
		require.NoError(t, err)

		accessCookie := cookies.NewAccessCookie(accessToken, atExp, suite.Config)
		refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
		req.AddCookie(&accessCookie)
		req.AddCookie(&refreshCookie)

		// Make the request
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify no new cookies are set (soft refresh doesn't update tokens when access token is valid)
		cookies := resp.Cookies()
		assert.Empty(t, cookies)
	})

	t.Run("calling a route with check & soft refresh, with [invalid access_token, valid refresh_token] returns a new [access_token]", func(t *testing.T) {
		t.Run("access_token is expired", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create a user
			user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
			require.NoError(t, err)
			user = suite.CreateUser(t, user, []string{"user"})

			// Create expired access token
			cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.Metadata)
			require.NoError(t, err)
			accessToken, atExp, err := jwtgen.Generate(cClaims.ToMap(), time.Now().Add(-2*time.Minute), 1, "TestApp", suite.Config.JWT.Secret)
			require.NoError(t, err)

			// Create valid refresh token
			refreshTokenEntity, _, err := entities.NewRefreshToken(user, "device-fingerprint", suite.Config)
			require.NoError(t, err)
			refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)

			// Create request with expired access token and valid refresh token
			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/me", nil)
			require.NoError(t, err)

			accessCookie := cookies.NewAccessCookie(accessToken, atExp, suite.Config)
			refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
			req.AddCookie(&accessCookie)
			req.AddCookie(&refreshCookie)

			// Make the request
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			// todo: must be fixed, actual error //////!\\\\\\
			//assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify new access token cookie is set
			cookies := resp.Cookies()
			var newAccessToken string
			//var newRefreshToken string
			for _, cookie := range cookies {
				if cookie.Name == "access_token" {
					newAccessToken = cookie.Value
					break
				}
				//if cookie.Name == "refresh_token" {
				//	newRefreshToken = cookie.Value
				//	break
				//}
			}
			assert.NotEmpty(t, newAccessToken)
			assert.NotEqual(t, accessToken, newAccessToken)

			// verify old refresh token is not there
			var count int64
			err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
			require.NoError(t, err)
			assert.Equal(t, int64(0), count)

			// todo: fix
			// verify new refresh token is there
			//var count2 int64
			//err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ? AND user_id = ?", newRefreshToken, user.ID).Count(&count2).Error
			//require.NoError(t, err)
			//assert.Equal(t, int64(1), count2)
		})

		// todo: fix like above

		//t.Run("access_token is malformed", func(t *testing.T) {
		//	suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
		//	defer suite.Teardown()
		//
		//	// Create a user
		//	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
		//	require.NoError(t, err)
		//	user = suite.CreateUser(t, user, []string{"user"})
		//
		//	// Create malformed access token
		//	malformedAccessToken := "invalid.jwt.token"
		//
		//	// Create valid refresh token
		//	refreshTokenEntity, _, err := entities.NewRefreshToken(user, "device-fingerprint", suite.Config)
		//	require.NoError(t, err)
		//	refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)
		//
		//	// Create request with malformed access token and valid refresh token
		//	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/me", nil)
		//	require.NoError(t, err)
		//
		//	accessCookie := cookies.NewAccessCookie(malformedAccessToken, time.Now().Add(time.Hour).Unix(), suite.Config)
		//	refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
		//	req.AddCookie(&accessCookie)
		//	req.AddCookie(&refreshCookie)
		//
		//	// Make the request
		//	resp, err := http.DefaultClient.Do(req)
		//	require.NoError(t, err)
		//	assert.Equal(t, http.StatusOK, resp.StatusCode)
		//
		//	// Verify new access token cookie is set
		//	cookies := resp.Cookies()
		//	var newAccessToken string
		//	for _, cookie := range cookies {
		//		if cookie.Name == "access_token" {
		//			newAccessToken = cookie.Value
		//			break
		//		}
		//	}
		//	assert.NotEmpty(t, newAccessToken)
		//	assert.NotEqual(t, malformedAccessToken, newAccessToken)
		//
		//	// verify old refresh token is not there
		//	var count int64
		//	err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
		//	require.NoError(t, err)
		//	assert.Equal(t, int64(0), count)
		//})
		//
		//t.Run("access_token is empty", func(t *testing.T) {
		//	suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
		//	defer suite.Teardown()
		//
		//	// Create a user
		//	user, err := entities.NewUser("testuser", "https://example.com/avatar.jpg", "test@example.com", "github")
		//	require.NoError(t, err)
		//	user = suite.CreateUser(t, user, []string{"user"})
		//
		//	// Create valid refresh token
		//	refreshTokenEntity, _, err := entities.NewRefreshToken(user, "device-fingerprint", suite.Config)
		//	require.NoError(t, err)
		//	refreshTokenEntity = suite.CreateRefreshToken(t, refreshTokenEntity)
		//
		//	// Create request with empty access token and valid refresh token
		//	req, err := http.NewRequest("GET", suite.Server.URL+"/auth/me", nil)
		//	require.NoError(t, err)
		//
		//	refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
		//	req.AddCookie(&refreshCookie)
		//
		//	// Make the request
		//	resp, err := http.DefaultClient.Do(req)
		//	require.NoError(t, err)
		//	assert.Equal(t, http.StatusOK, resp.StatusCode)
		//
		//	// Verify new access token cookie is set
		//	cookies := resp.Cookies()
		//	var newAccessToken string
		//	for _, cookie := range cookies {
		//		if cookie.Name == "access_token" {
		//			newAccessToken = cookie.Value
		//			break
		//		}
		//	}
		//	assert.NotEmpty(t, newAccessToken)
		//
		//	// verify old refresh token is still there
		//	var count int64
		//	err = suite.Db.Model(&entities.RefreshToken{}).Where("token = ?", refreshTokenEntity.Token).Count(&count).Error
		//	require.NoError(t, err)
		//	assert.Equal(t, int64(1), count)
		//})
	})

	t.Run("calling a route with check & soft refresh, with [invalid access_token, invalid refresh_token] returns 401", func(t *testing.T) {
		t.Run("refresh_token is expired", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
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
			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/me", nil)
			require.NoError(t, err)

			refreshCookie := cookies.NewRefreshCookie(refreshTokenEntity.Token, refreshTokenEntity.ExpiresAt.Unix(), suite.Config)
			req.AddCookie(&refreshCookie)

			// Make the request
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			// Verify we dont get cookies
			cookies := resp.Cookies()
			assert.Empty(t, cookies)

			// Verify error response
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			var errorResponse map[string]interface{}
			err = json.Unmarshal(body, &errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse["error"], apperrors.ErrRefreshTokenExpired.Error())
		})

		t.Run("refresh_token is malformed", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create request with malformed refresh token
			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/me", nil)
			require.NoError(t, err)

			refreshCookie := cookies.NewRefreshCookie("malformed_refresh_token", time.Now().Add(time.Hour).Unix(), suite.Config)
			req.AddCookie(&refreshCookie)

			// Make the request
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			// Verify we dont get cookies
			cookies := resp.Cookies()
			assert.Empty(t, cookies)

			// Verify error response
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			var errorResponse map[string]interface{}
			err = json.Unmarshal(body, &errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse["error"], apperrors.ErrGeneric.Error())
		})

		t.Run("refresh_token is empty", func(t *testing.T) {
			suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
			defer suite.Teardown()

			// Create request without refresh token
			req, err := http.NewRequest("GET", suite.Server.URL+"/auth/me", nil)
			require.NoError(t, err)

			// Make the request
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			// Verify we dont get cookies
			cookies := resp.Cookies()
			assert.Empty(t, cookies)

			// Verify error response
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			var errorResponse map[string]interface{}
			err = json.Unmarshal(body, &errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse["error"], apperrors.ErrGeneric.Error())
		})
	})
}

// t.Run("no refresh: calling a route with only a check (ex: /authorize) should never refresh tokens", func(t *testing.T) {
// 	// todo: test when /authorize is implemented
// })

// t.Run("rate limiting", func(t *testing.T) {
// 	// todo: implement and test rate limiting
// })

// t.Run("devices ids", func(t *testing.T) {
// 	// todo: implement and test different devices ids
// })
