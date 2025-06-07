package cookies

import (
	"net/http"
	"time"

	"aegix/internal/domain"
)

func NewAccessCookie(accessToken string, expiresAt int64, withDefaults bool, config domain.Config) http.Cookie {
	return newCookie("access_token", accessToken, expiresAt, withDefaults, config)
}

func NewRefreshCookie(refreshToken string, expiresAt int64, withDefaults bool, config domain.Config) http.Cookie {
	return newCookie("refresh_token", refreshToken, expiresAt, withDefaults, config)
}

func newCookie(name, value string, expiresAt int64, withDefaults bool, config domain.Config) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Domain:   config.App.URL,
		Value:    value,
		Expires:  time.Unix(expiresAt, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	if withDefaults {
		if config.Cookies.Domain != "" {
			cookie.Domain = config.Cookies.Domain
		}
		if config.Cookies.Path != "" {
			cookie.Path = config.Cookies.Path
		}
		if config.Cookies.Secure {
			cookie.Secure = true
		}
		if config.Cookies.HTTPOnly {
			cookie.HttpOnly = true
		}
		if config.Cookies.SameSite != "" {
			if config.Cookies.SameSite == "Lax" {
				cookie.SameSite = http.SameSiteLaxMode
			} else if config.Cookies.SameSite == "Strict" {
				cookie.SameSite = http.SameSiteStrictMode
			} else if config.Cookies.SameSite == "None" {
				cookie.SameSite = http.SameSiteNoneMode
			}
		}
	}
	return cookie
}
