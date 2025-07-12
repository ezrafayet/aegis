package integration

import (
	"aegis/integration-tests/testkit"
	"aegis/pkg/apperrors"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderCallback(t *testing.T) {
	t.Run("unhappy scenarios: invalid data / generic", func(t *testing.T) {
		t.Run("calling GET /provider/callback returns 403 if the provider is not enabled", func(t *testing.T) {
			config := testkit.GetBaseConfig()
			config.Auth.Providers.GitHub.Enabled = false
			suite := testkit.SetupTestSuite(t, config)
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
		t.Run("calling GET /provider/callback with invalid state gets rejected", func(t *testing.T) {})
		t.Run("calling GET /provider/callback with invalid code gets rejected", func(t *testing.T) {})
	})
	t.Run("unhappy scenarios: error page", func(t *testing.T) {
		t.Run("calling GET /provider/callback returns to error page if user declines auth", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is using another method", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is blocked", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is deleted", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is not an early user", func(t *testing.T) {})
	})
	t.Run("happy scenarios", func(t *testing.T) {
		t.Run("calling GET /provider/callback gives [access_token, refresh_token] if the user already exists", func(t *testing.T) {})
		t.Run("calling GET /provider/callback gives [access_token, refresh_token] and creates user if the user does not exist", func(t *testing.T) {})
		t.Run("calling GET /provider/callback cleans the state", func(t *testing.T) {})
		t.Run("calling GET /provider/callback redirects to the welcome page", func(t *testing.T) {})
	})
}
