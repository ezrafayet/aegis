package handlers

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/pkg/apperrors"
	"aegis/pkg/cookies"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HandlersInterface interface {
	GetSession(c echo.Context) error
	Logout(c echo.Context) error
	DoNothing(c echo.Context) error
}

type Handlers struct {
	Config  entities.Config
	Service primary.UseCasesInterface
}

var _ HandlersInterface = (*Handlers)(nil)

func NewHandlers(c entities.Config, s primary.UseCasesInterface) *Handlers {
	return &Handlers{
		Config:  c,
		Service: s,
	}
}

func (h Handlers) GetSession(c echo.Context) error {
	var accessToken string
	if cookie, err := c.Cookie("access_token"); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
	} else {
		accessToken = cookie.Value
	}
	session, err := h.Service.GetSession(accessToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrAccessTokenExpired) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenExpired.Error()})
		}
		if errors.Is(err, apperrors.ErrAccessTokenInvalid) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenInvalid.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	return c.JSON(http.StatusOK, session)
}

func (h Handlers) Logout(c echo.Context) error {
	var refreshToken string
	if cookie, err := c.Cookie("refresh_token"); err != nil {
		refreshToken = ""
	} else {
		refreshToken = cookie.Value
	}
	tokensPair, err := h.Service.Logout(refreshToken)
	if tokensPair != nil {
		accessCookie := cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), h.Config)
		refreshCookie := cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), h.Config)
		c.SetCookie(&accessCookie)
		c.SetCookie(&refreshCookie)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h Handlers) DoNothing(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
