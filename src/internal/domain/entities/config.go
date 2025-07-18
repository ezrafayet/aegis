package entities

type Config struct {
	App struct {
		// Name of the application (ex: "Aegis")
		Name string `json:"name"`
		// URL of the application (main domain, ex: https://aegis.example.com)
		URL string `json:"url"`
		// Allowed origins for the application (CORS) (ex: ["https://aegis.example.com", "http://localhost:5000])
		CorsAllowedOrigins []string `json:"cors_allowed_origins"`
		// New users need to be approved by an admin (ex: true)
		EarlyAdoptersOnly bool `json:"early_adopters_only"`
		// Redirect URL after successful login (ex: "https://aegis.example.com" or https://aegis.example.com/login-success)
		RedirectAfterSuccess string `json:"redirect_after_success"`
		// Redirect URL after login error (ex: "https://aegis.example.com/login-error")
		RedirectAfterError string `json:"redirect_after_error"`
		// API keys for the application (used for internal requests) (ex: ["1234567890"])
		InternalAPIKeys []string `json:"internal_api_keys"`
		// Port on which the service must run (ex: 5666)
		Port int `json:"port"`
	} `json:"app"`

	LoginPage struct {
		// If true, the login page will be enabled
		Enabled bool `json:"enabled"`
		// Full path to the login page (ex: "/login")
		FullPath string `json:"full_path"`
	} `json:"login_page"`

	ErrorPage struct {
		// If true, the error page will be enabled
		Enabled bool `json:"enabled"`
		// Full path to the error page (ex: "/login-error")
		FullPath string `json:"full_path"`
	} `json:"error_page"`

	DB struct {
		// DB connection string (ex: "postgres://user:password@localhost:5432/auth")
		PostgresURL string `json:"postgres_url"`
	} `json:"db"`

	JWT JWTConfig `json:"jwt"`

	Auth struct {
		// Providers configuration
		Providers struct {
			GitHub struct {
				Enabled      bool   `json:"enabled"`
				AppName      string `json:"app_name"`
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
			} `json:"github"`
			Discord struct {
				Enabled      bool   `json:"enabled"`
				AppName      string `json:"app_name"`
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
			} `json:"discord"`
		} `json:"providers"`
	} `json:"auth"`

	Cookies struct {
		Domain   string `json:"domain"`
		Secure   bool   `json:"secure"`
		HTTPOnly bool   `json:"http_only"`
		// SameSite cookie attribute: 1 = default, 2 = lax, 3 = strict, 4 = none
		SameSite int    `json:"same_site"`
		Path     string `json:"path"`
	} `json:"cookies"`

	User struct {
		// Roles for a user. Mandatory roles are: "user" and "platform_admin"
		Roles []string `json:"roles"`
	} `json:"user"`
}

type JWTConfig struct {
	// Secret key for the JWT
	Secret string `json:"secret"`
	// Access token expiration time in minutes
	AccessTokenExpirationMin int `json:"access_token_expiration_minutes"`
	// Refresh token expiration time in days
	RefreshTokenExpirationDays int `json:"refresh_token_expiration_days"`
}
