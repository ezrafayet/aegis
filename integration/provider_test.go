package integration

import (
	"aegis/integration/integration_testkit"
	"aegis/internal/domain/entities"
	"aegis/pkg/apperrors"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	t.Run("calling GET /provider returns 403 if the provider is not enabled", func(t *testing.T) {
		config := integration_testkit.GetBaseConfig()
		config.Auth.Providers.GitHub.Enabled = false
		suite := integration_testkit.SetupTestSuite(t, config)
		defer suite.Teardown()

		resp, err := http.Get(suite.Server.URL + "/auth/github")
		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var errorResponse map[string]interface{}
		err = json.Unmarshal(body, &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], apperrors.ErrAuthMethodNotEnabled.Error())
	})

	t.Run("calling GET /provider returns the url for client to redirect to, and saves a state", func(t *testing.T) {
		config := integration_testkit.GetBaseConfig()
		config.Auth.Providers.GitHub.Enabled = true
		suite := integration_testkit.SetupTestSuite(t, config)
		defer suite.Teardown()

		resp, err := http.Get(suite.Server.URL + "/auth/github")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)

		// Check that the response contains a URL
		assert.Contains(t, response, "redirect_url")
		url, ok := response["redirect_url"].(string)
		assert.True(t, ok)
		assert.NotEmpty(t, url)
		assert.Contains(t, url, "github.com")
		assert.Contains(t, url, "client_id="+config.Auth.Providers.GitHub.ClientID)

		// Check that a state parameter is present in the URL
		assert.Contains(t, url, "state=")

		// Verify that a state was saved in the database
		var stateCount int64
		err = suite.Db.Model(&entities.State{}).Count(&stateCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), stateCount)
	})

	t.Run("rate limiting", func(t *testing.T) {
		// todo: implement and test rate limiting
	})
}
