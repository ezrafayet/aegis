package discord

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OAuthDiscordRepository struct {
	Name   string
	Config entities.Config
}

var _ secondary.OAuthProviderInterface = OAuthDiscordRepository{}

func NewOAuthDiscordRepository(c entities.Config) OAuthDiscordRepository {
	return OAuthDiscordRepository{
		Name:   "discord",
		Config: c,
	}
}

type discordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type discordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
}

func (p OAuthDiscordRepository) IsEnabled() (bool, error) {
	fmt.Println("Config hit")
	fmt.Println("IsEnabled", p.Config.Auth.Providers.Discord.Enabled)
	return p.Config.Auth.Providers.Discord.Enabled, nil
}

func (p OAuthDiscordRepository) GetOauthRedirectURL(redirectUrl, state string) (string, error) {
	fmt.Println("GetOauthRedirectURL", redirectUrl, state, p.Config.Auth.Providers.Discord.ClientID)
	return fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20email&state=%s",
		p.Config.Auth.Providers.Discord.ClientID,
		redirectUrl,
		state,
	), nil
}

func (p OAuthDiscordRepository) GetName() string {
	fmt.Println("GetName", p.Name)
	return p.Name
}

func (p OAuthDiscordRepository) ExchangeCodeForUserInfos(code, state, redirectUri string) (*entities.UserInfos, error) {
	// Step 1: get access token
	data := map[string]string{
		"client_id":     p.Config.Auth.Providers.Discord.ClientID,
		"client_secret": p.Config.Auth.Providers.Discord.ClientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectUri,
	}
	fmt.Println(data)

	// Convert data to form-encoded format
	formData := ""
	for key, value := range data {
		if formData != "" {
			formData += "&"
		}
		formData += key + "=" + value
	}

	req1, _ := http.NewRequest("POST", "https://discord.com/api/oauth2/token", bytes.NewBufferString(formData))
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp1.Body.Close()
	var tokenResponse discordTokenResponse
	if err := json.NewDecoder(resp1.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode access token: %w", err)
	}

	// Step 2: get user infos
	req2, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	req2.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp2.Body.Close()
	var user discordUser
	if err := json.NewDecoder(resp2.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Construct avatar URL if avatar exists
	var avatarURL string
	if user.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.ID, user.Avatar)
	}

	// Construct display name
	displayName := user.Username
	if user.Discriminator != "0" {
		displayName = fmt.Sprintf("%s#%s", user.Username, user.Discriminator)
	}

	return &entities.UserInfos{
		Name:   displayName,
		Email:  user.Email,
		Avatar: avatarURL,
	}, nil
}
