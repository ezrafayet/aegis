package integration

import (
	"aegis/internal/domain/entities"
	"aegis/pkg/cookies"
	"aegis/pkg/jwtgen"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMe(t *testing.T) {
	t.Run("asking /me without a session returns 401", func(t *testing.T) {
		suite := setupTestSuite(t)
		defer suite.teardown()
		resp, err := http.Get(suite.server.URL + "/auth/me")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
	t.Run("asking /me with a session returns 200 and the session", func(t *testing.T) {
		suite := setupTestSuite(t)
		defer suite.teardown()
		user, err := entities.NewUser("cloude", "https://example.com/avatar.jpg", "cloude@example.com", "github")
		require.NoError(t, err)
		user = suite.CreateUser(t, user, []string{"user"})
		require.NoError(t, err)
		cClaims, err := entities.NewCustomClaimsFromValues(user.ID, false, user.Roles, user.Metadata)
		require.NoError(t, err)
		accessToken, atExp, err := jwtgen.Generate(cClaims.ToMap(), time.Now(), 10, "MyApp", suite.config.JWT.Secret)
		require.NoError(t, err)
		atCookie := cookies.NewAccessCookie(accessToken, atExp, suite.config)
		req, err := http.NewRequest("GET", suite.server.URL+"/auth/me", nil)
		require.NoError(t, err)
		req.AddCookie(&atCookie)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var sessionResponse entities.Session
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&sessionResponse)
		require.NoError(t, err)
		assert.Equal(t, user.ID, sessionResponse.UserID)
		assert.Equal(t, "user", sessionResponse.Roles)
		assert.Equal(t, "{}", sessionResponse.Metadata)
		assert.Equal(t, user.EarlyAdopter, sessionResponse.EarlyAdopter)
	})
}
