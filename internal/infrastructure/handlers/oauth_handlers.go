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

type OAuthHandlersInterface interface {
	GetAuthURL(c echo.Context) error
	ExchangeCode(c echo.Context) error
}

type OAuthHandlers struct {
	Config  entities.Config
	Service primary.OAuthUseCasesInterface
}

var _ OAuthHandlersInterface = OAuthHandlers{}

func NewOAuthHandlers(c entities.Config, s primary.OAuthUseCasesInterface) OAuthHandlers {
	return OAuthHandlers{
		Config:  c,
		Service: s,
	}
}

func (h OAuthHandlers) GetAuthURL(c echo.Context) error {
	redirectUrl, err := h.Service.GetAuthURL(c.QueryParam("redirect_uri"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an error occurred"})
	}
	return c.JSON(http.StatusOK, map[string]string{"redirect_url": redirectUrl})
}

func (h OAuthHandlers) ExchangeCode(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	error := c.QueryParam("error")
	if error != "" {
		return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterError)
	}
	tokensPair, err := h.Service.ExchangeCode(code, state)
	if err != nil {
		if errors.Is(err, apperrors.ErrWrongAuthMethod) {
			return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterError)
		}
		if errors.Is(err, apperrors.ErrEarlyAdoptersOnly) {
			return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterError)
		}
		if errors.Is(err, apperrors.ErrUserBlocked) {
			return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterError)
		}
		if errors.Is(err, apperrors.ErrUserDeleted) {
			return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterError)
		}
		return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterError)
	}
	if tokensPair != nil {
		accessCookie := cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), h.Config)
		refreshCookie := cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), h.Config)

		c.SetCookie(&accessCookie)
		c.SetCookie(&refreshCookie)
	}
	return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterSuccess)
}
