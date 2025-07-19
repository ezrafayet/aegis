package integration_test_cases

import (
	"aegis/integration/integration_testkit"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func LoginError_EnabledReturns200AndShowsErrorPage(t *testing.T) {
	config := integration_testkit.GetBaseConfig()
	// Enable error page
	config.ErrorPage.Enabled = true
	config.ErrorPage.FullPath = "/auth/login-error"

	suite := integration_testkit.SetupTestSuite(t, config)
	defer suite.Teardown()

	resp, err := http.Get(suite.Server.URL + config.ErrorPage.FullPath)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "<!DOCTYPE html>")
	assert.Contains(t, string(body), "TestApp")
}

func LoginError_DisabledReturns404(t *testing.T) {
	config := integration_testkit.GetBaseConfig()
	// Disable error page
	config.ErrorPage.Enabled = false
	config.ErrorPage.FullPath = "/auth/login-error"

	suite := integration_testkit.SetupTestSuite(t, config)
	defer suite.Teardown()

	resp, err := http.Get(suite.Server.URL + config.ErrorPage.FullPath)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
