package handlers

import (
	"net/http"
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/primary_ports"
	"othnx/pkg/apperrors"
	"othnx/pkg/cookies"

	"github.com/labstack/echo/v4"
)

type HandlersInterface interface {
	GetSession(c echo.Context) error
	Logout(c echo.Context) error
	DoNothing(c echo.Context) error
}

type Handlers struct {
	Config  entities.Config
	Service primaryports.UseCasesInterface
}

var _ HandlersInterface = &Handlers{}

func NewHandlers(c entities.Config, s primaryports.UseCasesInterface) Handlers {
	return Handlers{
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
		if err.Error() == apperrors.ErrAccessTokenExpired.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenExpired.Error()})
		}
		if err.Error() == apperrors.ErrAccessTokenInvalid.Error() {
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
		cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), h.Config)
		cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), h.Config)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h Handlers) DoNothing(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
