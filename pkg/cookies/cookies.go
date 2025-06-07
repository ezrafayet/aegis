package cookies

import (
	"net/http"
	"time"

	"aegix/internal/domain"
)

func NewAccessCookie(accessToken string, expiresAt int64, config domain.Config) http.Cookie {
	return newCookie("access_token", accessToken, expiresAt, config)
}

func NewRefreshCookie(refreshToken string, expiresAt int64, config domain.Config) http.Cookie {
	return newCookie("refresh_token", refreshToken, expiresAt, config)
}

func newCookie(name, value string, expiresAt int64, config domain.Config) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Domain:   config.Cookies.Domain,
		Value:    value,
		Expires:  time.Unix(expiresAt, 0),
		HttpOnly: config.Cookies.HTTPOnly,
		Secure:   config.Cookies.Secure,
		SameSite: http.SameSite(config.Cookies.SameSite),
		Path:     config.Cookies.Path,
	}
	return cookie
}
