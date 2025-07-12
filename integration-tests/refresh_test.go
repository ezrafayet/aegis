package integration

import "testing"

func TestRefresh(t *testing.T) {
	t.Run("calling GET /refresh with [valid access_token, valid refresh_token] returns a new [access_token, refresh_token]", func(t *testing.T) {})
	t.Run("calling GET /refresh with [invalid access_token, valid refresh_token] returns a new [access_token, refresh_token]", func(t *testing.T) {
		t.Run("access_token is expired", func(t *testing.T) {})
		t.Run("access_token is malformed", func(t *testing.T) {})
		t.Run("access_token is empty", func(t *testing.T) {})
	})
	t.Run("calling GET /refresh with [invalid access_token, invalid refresh_token] returns 401", func(t *testing.T) {
		t.Run("refresh_token is expired", func(t *testing.T) {})
		t.Run("refresh_token is malformed", func(t *testing.T) {})
		t.Run("refresh_token is empty", func(t *testing.T) {})
	})
	t.Run("middleware testing", func(t *testing.T) {
		t.Run("calling a route with check & soft refresh, with [valid access_token, valid refresh_token] does not update the tokens", func(t *testing.T) {})
		t.Run("calling a route with check & soft refresh, with [invalid access_token, valid refresh_token] returns a new [access_token]", func(t *testing.T) {
			t.Run("access_token is expired", func(t *testing.T) {})
			t.Run("access_token is malformed", func(t *testing.T) {})
			t.Run("access_token is empty", func(t *testing.T) {})
		})
		t.Run("calling a route with check & soft refresh, with [invalid access_token, invalid refresh_token] returns 401", func(t *testing.T) {
			t.Run("refresh_token is expired", func(t *testing.T) {})
			t.Run("refresh_token is malformed", func(t *testing.T) {})
			t.Run("refresh_token is empty", func(t *testing.T) {})
		})
	})
}
