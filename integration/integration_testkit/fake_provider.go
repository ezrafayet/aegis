package integration_testkit

import (
	"aegis/pkg/plugins/providers"
	"fmt"
)

type FakeOAuthProvider struct {
	Name         string
	Enabled      bool
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var _ providers.OAuthProviderInterface = (*FakeOAuthProvider)(nil)

func NewFakeOAuthProvider(name string, enabled bool, clientID, clientSecret, redirectURL string) *FakeOAuthProvider {
	return &FakeOAuthProvider{
		Name:         name,
		Enabled:      enabled,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

func (p *FakeOAuthProvider) IsEnabled() bool {
	return p.Enabled
}

func (p *FakeOAuthProvider) GetName() string {
	return p.Name
}

func (p *FakeOAuthProvider) GetOauthRedirectURL(state string) string {
	return fmt.Sprintf(
		"https://%s.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s",
		p.Name,
		p.ClientID,
		state,
	)
}

func (p *FakeOAuthProvider) ExchangeCodeForUserInfos(code, state string) (*providers.UserInfos, error) {
	// Simple fake implementation that returns consistent test data
	// No HTTP calls needed!

	// Check if the code is valid (for testing different scenarios)
	switch code {
	case "accepted_code":
		// Return valid user info
		return &providers.UserInfos{
			Name:   "testuser",
			Email:  "test@example.com",
			Avatar: "https://example.com/avatar.jpg",
		}, nil
	case "rejected_code":
		return nil, fmt.Errorf("invalid_grant: The authorization code is invalid or has expired")
	case "declined_code":
		return nil, fmt.Errorf("access_denied: User declined to authorize the application")
	default:
		return nil, fmt.Errorf("unknown_error: Unexpected code")
	}
}
