package domain

type Config struct {
	App struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		LogLevel string `json:"log_level"`
		APIKeys  []string `json:"api_keys"`
	} `json:"app"`

	DB struct {
		PostgresURL string `json:"postgres_url"`
	} `json:"db"`

	JWT struct {
		Secret                     string `json:"secret"`
		AccessTokenExpirationMin   int    `json:"access_token_expiration_minutes"`
		RefreshTokenExpirationDays int    `json:"refresh_token_expiration_days"`
	} `json:"jwt"`

	Auth struct {
		Providers struct {
			GitHub struct {
				Enabled      bool   `json:"enabled"`
				AppName      string `json:"app_name"`
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
			} `json:"github"`
		} `json:"providers"`

		RedirectURLs struct {
			Success []string `json:"success"`
			Failure []string `json:"failure"`
		} `json:"redirect_urls"`

		AllowedOrigins []string `json:"allowed_origins"`
	} `json:"auth"`

	Cookie struct {
		Domain   string `json:"domain"`
		Secure   bool   `json:"secure"`
		HTTPOnly bool   `json:"http_only"`
		SameSite string `json:"same_site"`
		MaxAge   int    `json:"max_age"`
	} `json:"cookie"`

	User struct {
		Roles    []string                      `json:"roles"`
		Metadata map[string]UserMetadataConfig `json:"metadata"`
	} `json:"user"`
}

type UserMetadataConfig struct {
	Type    string   `json:"type"`
	Default string   `json:"default"`
	Enum    []string `json:"enum"`
}
