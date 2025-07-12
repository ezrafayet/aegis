package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	t.Run("calling GET /health returns 200", func(t *testing.T) {
		suite := setupTestSuite(t)
		defer suite.teardown()
		resp, err := http.Get(suite.server.URL + "/auth/health")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
