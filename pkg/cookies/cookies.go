package cookies

import (
	"net/http"
	"time"

	"aegix/internal/domain"
)

func NewAccessCookie(token string, expiresAt int64, config domain.Config) http.Cookie {
	return newCookie("access_token", token, expiresAt, config)
}

func NewRefreshCookie(token string, expiresAt int64, config domain.Config) http.Cookie {
	return newCookie("refresh_token", token, expiresAt, config)
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
