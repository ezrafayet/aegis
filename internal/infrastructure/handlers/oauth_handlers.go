package handlers

import (
	"errors"
	"net/http"
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/primary"
	"othnx/pkg/apperrors"
	"othnx/pkg/cookies"

	"github.com/labstack/echo/v4"
)

// some factory

type OAuthHandlersInterface interface {
	GetAuthURL(c echo.Context) error
	ExchangeCode(c echo.Context) error
}

type OAuthGithubHandlers struct {
	Config  entities.Config
	Service primary.OAuthUseCasesInterface
}

var _ OAuthHandlersInterface = OAuthGithubHandlers{}

func NewOAuthGithubHandlers(c entities.Config, s primary.OAuthUseCasesInterface) OAuthGithubHandlers {
	return OAuthGithubHandlers{
		Config:  c,
		Service: s,
	}
}

func (h OAuthGithubHandlers) GetAuthURL(c echo.Context) error {
	redirectUrl, err := h.Service.GetAuthURL(c.QueryParam("redirect_uri"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an error occurred"})
	}
	return c.JSON(http.StatusOK, map[string]string{"redirect_url": redirectUrl})
}

func (h OAuthGithubHandlers) ExchangeCode(c echo.Context) error {
	type ExchangeCodeRequest struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
	var body ExchangeCodeRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	tokensPair, err := h.Service.ExchangeCode(body.Code, body.State)
	if err != nil {
		if errors.Is(err, apperrors.ErrEarlyAdoptersOnly) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrEarlyAdoptersOnly.Error()})
		}
		if errors.Is(err, apperrors.ErrUserBlocked) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserBlocked.Error()})
		}
		if errors.Is(err, apperrors.ErrUserDeleted) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserDeleted.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	if tokensPair != nil {
		cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), h.Config)
		cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), h.Config)
	}
	return c.NoContent(http.StatusOK)
}
