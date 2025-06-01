package github

import (
	"aegix/internal/domain"
	"aegix/internal/providers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OAuthGithubProvider struct {
	Config domain.Config
}

var _ providers.OAuthProvider = OAuthGithubProvider{}

func NewOAuthGithubProvider(c domain.Config) OAuthGithubProvider {
	return OAuthGithubProvider{
		Config: c,
	}
}

type gitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type gitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (p OAuthGithubProvider) GetUserInfos(code, state, redirectUri string) (*providers.OAuthUser, error) {
	fmt.Println("hit1")
	// Step 1: get access token
	data := map[string]string{
		"client_id":     p.Config.Auth.Providers.GitHub.ClientID,
		"client_secret": p.Config.Auth.Providers.GitHub.ClientSecret,
		"code":          code,
		"state":         state,
	}
	fmt.Println(data)
	if redirectUri != "" {
		data["redirect_uri"] = redirectUri
	}
	body1, _ := json.Marshal(data)
	req1, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(body1))
	req1.Header.Set("Accept", "application/json")
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp1.Body.Close()
	var tokenResponse gitHubTokenResponse
	if err := json.NewDecoder(resp1.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode access token: %w", err)
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> tokenResponse", tokenResponse)
	// Step 2: get user infos
	// Step 3: get user emails
	req3, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}
	req3.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)
	// req3.Header.Set("User-Agent", p.Config.Auth.Providers.GitHub.AppName)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}
	defer resp3.Body.Close()
	var emails []gitHubEmail
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> resp3", resp3)
	if err := json.NewDecoder(resp3.Body).Decode(&emails); err != nil {
		return nil, fmt.Errorf("failed to decode user emails: %w", err)
	}
	var em string
	for _, email := range emails {
		if email.Primary && email.Verified {
			em = email.Email
		}
	}
	return &providers.OAuthUser{
		Name: "",
		Email: em,
		Avatar: "",
	}, nil
}
