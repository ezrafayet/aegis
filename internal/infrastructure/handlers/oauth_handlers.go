package handlers

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/pkg/apperrors"
	"aegis/pkg/cookies"
	"aegis/pkg/urlbuilder"
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

var _ OAuthHandlersInterface = (*OAuthHandlers)(nil)

func NewOAuthHandlers(c entities.Config, s primary.OAuthUseCasesInterface) *OAuthHandlers {
	return &OAuthHandlers{
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
		redirectURL, err := urlbuilder.Build(h.Config.App.RedirectAfterError, "", map[string]string{"error": error})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an error occurred"})
		}
		return c.Redirect(http.StatusFound, redirectURL)
	}
	tokensPair, err := h.Service.ExchangeCode(code, state)
	if err != nil {
		var errorType string
		if errors.Is(err, apperrors.ErrWrongAuthMethod) {
			errorType = "wrong_auth_method"
		} else if errors.Is(err, apperrors.ErrEarlyAdoptersOnly) {
			errorType = "early_adopters_only"
		} else if errors.Is(err, apperrors.ErrUserBlocked) {
			errorType = "user_blocked"
		} else if errors.Is(err, apperrors.ErrUserDeleted) {
			errorType = "user_deleted"
		} else {
			errorType = "unknown_error"
		}
		redirectURL, err := urlbuilder.Build(h.Config.App.RedirectAfterError, "", map[string]string{"error": errorType})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an error occurred"})
		}
		return c.Redirect(http.StatusFound, redirectURL)
	}
	if tokensPair != nil {
		accessCookie := cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), h.Config)
		refreshCookie := cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), h.Config)

		c.SetCookie(&accessCookie)
		c.SetCookie(&refreshCookie)
	}
	return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterSuccess)
}
