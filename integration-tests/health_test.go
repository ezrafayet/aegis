package integration

import (
	"net/http"
	"testing"

	"aegis/integration-tests/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	t.Run("calling GET /health returns 200", func(t *testing.T) {
		suite := testkit.SetupTestSuite(t, testkit.GetBaseConfig())
		defer suite.Teardown()
		resp, err := http.Get(suite.Server.URL + "/auth/health")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
