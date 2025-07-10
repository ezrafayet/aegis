package discord

import (
	"aegis/pkg/plugins/providers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OAuthDiscordRepository providers.OAuthRepository

var _ providers.OAuthProviderInterface = OAuthDiscordRepository{}

func NewOAuthDiscordRepository(enabled bool, clientID, clientSecret, redirectURL string) OAuthDiscordRepository {
	return OAuthDiscordRepository{
		Name:   "discord",
		Enabled: enabled,
		ClientID: clientID,
		ClientSecret: clientSecret,
		RedirectURL: redirectURL,
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

func (p OAuthDiscordRepository) IsEnabled() bool {
	return p.Enabled
}

func (p OAuthDiscordRepository) GetOauthRedirectURL(state string) string {
	return fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20email&state=%s",
		p.ClientID,
		p.RedirectURL,
		state,
	)
}

func (p OAuthDiscordRepository) GetName() string {
	fmt.Println("GetName", p.Name)
	return p.Name
}

func (p OAuthDiscordRepository) ExchangeCodeForUserInfos(code, state string) (*providers.UserInfos, error) {
	data := map[string]string{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  p.RedirectURL,
	}

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

	// Check if the response is successful
	if resp1.StatusCode != http.StatusOK {
		// Read error response body
		var errorBody bytes.Buffer
		errorBody.ReadFrom(resp1.Body)
		return nil, fmt.Errorf("discord token exchange failed with status %d: %s", resp1.StatusCode, errorBody.String())
	}

	var tokenResponse discordTokenResponse
	if err := json.NewDecoder(resp1.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode access token: %w", err)
	}

	// Check if we got a valid access token
	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("no access token received from Discord")
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
	// Check if the user info response is successful
	if resp2.StatusCode != http.StatusOK {
		// Read error response body
		var errorBody bytes.Buffer
		errorBody.ReadFrom(resp2.Body)
		return nil, fmt.Errorf("discord user info failed with status %d: %s", resp2.StatusCode, errorBody.String())
	}

	var user discordUser
	if err := json.NewDecoder(resp2.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Check if we got valid user data
	if user.ID == "" {
		fmt.Println("ERROR: No user ID received")
		return nil, fmt.Errorf("no user ID received from Discord")
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

	result := &providers.UserInfos{
		Name:   displayName,
		Email:  user.Email,
		Avatar: avatarURL,
	}

	return result, nil
}
