package integration_test_cases

import (
	"aegis/integration/integration_testkit"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Login_Returns200(t *testing.T) {
	config := integration_testkit.GetBaseConfig()
	config.LoginPage.Enabled = true

	suite := integration_testkit.SetupTestSuite(t, config)
	defer suite.Teardown()

	resp, err := http.Get(suite.Server.URL + config.LoginPage.FullPath)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
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
