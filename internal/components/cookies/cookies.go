package cookies

import (
	"aegix/internal/domain"
	"net/http"
	"time"
)

type CookieBuilderMethods interface {
	NewAccessCookie(accessToken string, expiresAt int64, withDefaults bool) http.Cookie
	NewRefreshCookie(refreshToken string, expiresAt int64, withDefaults bool) http.Cookie
}

type CookieBuilder struct {
	Config domain.Config
}

var _ CookieBuilderMethods = &CookieBuilder{}

func NewCookieBuilder(config domain.Config) CookieBuilderMethods {
	return &CookieBuilder{
		Config: config,
	}
}

func (b *CookieBuilder) NewAccessCookie(accessToken string, expiresAt int64, withDefaults bool) http.Cookie {
	return b.newCookie("access_token", accessToken, expiresAt, withDefaults)
}

func (b *CookieBuilder) NewRefreshCookie(refreshToken string, expiresAt int64, withDefaults bool) http.Cookie {
	return b.newCookie("refresh_token", refreshToken, expiresAt, withDefaults)
}

func (b *CookieBuilder) newCookie(name, value string, expiresAt int64, withDefaults bool) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Domain:   b.Config.App.URL,
		Value:    value,
		Expires:  time.Unix(expiresAt, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	if withDefaults {
		if b.Config.Cookie.Domain != "" {
			cookie.Domain = b.Config.Cookie.Domain
		}
		if b.Config.Cookie.Path != "" {
			cookie.Path = b.Config.Cookie.Path
		}
		if b.Config.Cookie.Secure {
			cookie.Secure = true
		}
		if b.Config.Cookie.HTTPOnly {
			cookie.HttpOnly = true
		}
		if b.Config.Cookie.SameSite != "" {
			if b.Config.Cookie.SameSite == "Lax" {
				cookie.SameSite = http.SameSiteLaxMode
			} else if b.Config.Cookie.SameSite == "Strict" {
				cookie.SameSite = http.SameSiteStrictMode
			} else if b.Config.Cookie.SameSite == "None" {
				cookie.SameSite = http.SameSiteNoneMode
			}
		}
	}
	return cookie
}
