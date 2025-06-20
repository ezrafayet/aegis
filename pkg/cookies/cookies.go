package cookies

import (
	"net/http"
	"time"

	"othnx/internal/domain"
)

func NewAccessCookie(token string, expiresAt int64, config domain.Config) http.Cookie {
	return newCookie("access_token", token, expiresAt, config)
}

func NewRefreshCookie(token string, expiresAt int64, config domain.Config) http.Cookie {
	return newCookie("refresh_token", token, expiresAt, config)
}

func NewAccessCookieZero(config domain.Config) http.Cookie {
	return newCookie("access_token", "", 0, config)
}

func NewRefreshCookieZero(config domain.Config) http.Cookie {
	return newCookie("refresh_token", "", 0, config)
}

func newCookie(name, token string, expiresAt int64, config domain.Config) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Domain:   config.Cookies.Domain,
		Value:    token,
		Expires:  time.Unix(expiresAt, 0),
		HttpOnly: config.Cookies.HTTPOnly,
		Secure:   config.Cookies.Secure,
		SameSite: http.SameSite(config.Cookies.SameSite),
		Path:     config.Cookies.Path,
	}
	return cookie
}

func IsZeroCookie(cookie http.Cookie) bool {
	return cookie.Expires.IsZero() || cookie.Value == ""
}
