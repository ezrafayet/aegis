package integration

import "testing"

func TestProvider(t *testing.T) {
	t.Run("calling GET /provider returns an error if the provider is not enabled", func(t *testing.T) {})

	t.Run("calling GET /provider returns the url for client to redirect to, and saves a state", func(t *testing.T) {})
}
