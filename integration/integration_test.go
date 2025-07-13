package integration

import (
	"aegis/integration/integration_test_cases"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Run("calling GET /health", func(t *testing.T) {
		t.Run("returns 200", integration_test_cases.Health_200)
	})
	t.Run("calling GET /logout", func(t *testing.T) {
		t.Run("sets zero cookies", integration_test_cases.Logout_SetsZeroCookies)
		t.Run("without a refresh_token does not break", integration_test_cases.Logout_WithoutRefreshTokenDoesNotBreak)
		t.Run("with a refresh_token, refresh_token gets deleted", integration_test_cases.Logout_DeletesRefreshToken)
		t.Run("rate limiting", func(t *testing.T) { /*todo*/ })
	})
	t.Run("calling GET /me", func(t *testing.T) {
		t.Run("without a session returns 401", integration_test_cases.Me_WithoutSessionReturns401)
		t.Run("with a session returns 200 and the session", integration_test_cases.Me_WithSessionReturns200)
		t.Run("rate limiting", func(t *testing.T) { /*todo*/ })
	})
	t.Run("calling GET /provider/callback", func(t *testing.T) {
		t.Run("unhappy scenarios: generic cases", func(t *testing.T) {
			t.Run("returns 403 if the provider is not enabled", integration_test_cases.ProviderCallback_NotEnabledReturns403)
		})
		t.Run("unhappy scenarios: cases that must redirect to error page", func(t *testing.T) {
			t.Run("redirects to error page if state is invalid", integration_test_cases.ProviderCallback_MustRedirectToErrorPage_InvalidState)
			t.Run("redirects to error page if code is invalid", integration_test_cases.ProviderCallback_MustRedirectToErrorPage_InvalidCode)
			t.Run("redirects to error page if user declines auth", integration_test_cases.ProviderCallback_MustRedirectToErrorPage_UserDeclinesAuth)
			t.Run("redirects to error page if user is using another method", integration_test_cases.ProviderCallback_MustRedirectToErrorPage_UserUsingAnotherMethod)
			t.Run("redirects to error page if user is blocked", integration_test_cases.ProviderCallback_MustRedirectToErrorPage_UserBlocked)
			t.Run("redirects to error page if user is deleted", integration_test_cases.ProviderCallback_MustRedirectToErrorPage_UserDeleted)
			t.Run("redirects to error page if user is not an early adopter", integration_test_cases.ProviderCallback_MustRedirectToErrorPage_UserNotAnEarlyAdopter)
		})
		t.Run("happy scenarios", func(t *testing.T) {
			t.Run("gives [access_token, refresh_token] if the user already exists", integration_test_cases.ProviderCallback_Success_UserExists)
			t.Run("gives [access_token, refresh_token] and creates userif the user does not exist", integration_test_cases.ProviderCallback_Success_UserDoesNotExist)
			t.Run("cleans the state", integration_test_cases.ProviderCallback_Success_CleansState)
			t.Run("redirects to the welcome page", integration_test_cases.ProviderCallback_Success_RedirectsToWelcomePage)
		})
		t.Run("rate limiting", func(t *testing.T) { /*todo*/ })
		t.Run("devices ids", func(t *testing.T) { /*todo*/ })
	})
	t.Run("calling GET /provider", func(t *testing.T) {
		t.Run("returns 403 if the provider is not enabled", integration_test_cases.Provider_NotEnabledReturns403)
		t.Run("returns 200 if the provider is enabled", integration_test_cases.Provider_EnabledReturnsUrlAndState)
		t.Run("rate limiting", func(t *testing.T) { /*todo*/ })
	})
	t.Run("calling GET /refresh", func(t *testing.T) {
		t.Run("hard refresh (always generate new tokens, ex: /refresh)", func(t *testing.T) {
			t.Run("must not refresh the user (1)", func(t *testing.T) {
				t.Run("if user does not exist", func(t *testing.T) { /*todo*/ })
				t.Run("if user is deleted", func(t *testing.T) { /*todo*/ })
				t.Run("if user is blocked", func(t *testing.T) { /*todo*/ })
				t.Run("if user is not an early adopter", func(t *testing.T) { /*todo*/ })
			})
			t.Run("must refresh the user", func(t *testing.T) {
				t.Run("if access_token and refresh_token are valid", integration_test_cases.HardRefresh_MustRefresh_ValidTokens)
				t.Run("if no access_token but refresh_token is valid", integration_test_cases.HardRefresh_MustRefresh_EmptyAT)
				t.Run("if expired access_token and refresh_token is valid", integration_test_cases.HardRefresh_MustRefresh_ExpiredAT)
				t.Run("if malformed access_token and refresh_token is valid", integration_test_cases.HardRefresh_MustRefresh_MalformedAT)
			})
			t.Run("must not refresh the user (2)", func(t *testing.T) {
				// todo: could be improved to cover the 9 cases
				t.Run("if no refresh_token", integration_test_cases.HardRefresh_MustNotRefresh_EmptyRT)
				t.Run("if expired refresh_token", integration_test_cases.HardRefresh_MustNotRefresh_ExpiredRT)
				t.Run("if malformed refresh_token", integration_test_cases.HardRefresh_MustNotRefresh_MalformedRT)
			})
		})
		t.Run("soft refresh (generate new tokens if this is needed only, ex: /me)", func(t *testing.T) {
			t.Run("must not refresh the user (1)", func(t *testing.T) {
				t.Run("if user does not exist", func(t *testing.T) { /*todo*/ })
				t.Run("if user is deleted", func(t *testing.T) { /*todo*/ })
				t.Run("if user is blocked", func(t *testing.T) { /*todo*/ })
				t.Run("if user is not an early adopter", func(t *testing.T) { /*todo*/ })
			})
			t.Run("must refresh the user", func(t *testing.T) {
				t.Run("if no access_token but refresh_token is valid", integration_test_cases.SoftRefresh_MustRefresh_EmptyAT)
				t.Run("if expired access_token and refresh_token is valid", integration_test_cases.SoftRefresh_MustRefresh_ExpiredAT)
				t.Run("if malformed access_token and refresh_token is valid", integration_test_cases.SoftRefresh_MustRefresh_MalformedAT)
			})
			t.Run("must not refresh the user (2)", func(t *testing.T) {
				// todo: could be improved to cover the 9 cases
				t.Run("if no refresh_token", integration_test_cases.SoftRefresh_MustNotRefresh_EmptyRT)
				t.Run("if expired refresh_token", integration_test_cases.SoftRefresh_MustNotRefresh_ExpiredRT)
				t.Run("if malformed refresh_token", integration_test_cases.SoftRefresh_MustNotRefresh_MalformedRT)
			})

			t.Run("must not refresh the user (3)", func(t *testing.T) {
				t.Run("if access_token is valid", integration_test_cases.SoftRefresh_MustNotRefresh_ValidAT)
			})
		})
		t.Run("rate limiting", func(t *testing.T) { /*todo*/ })
	})
}
