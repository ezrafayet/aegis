package handlers

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/pkg/apperrors"
	"aegis/pkg/cookies"
	"embed"
	"errors"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html
var templates embed.FS

type HandlersInterface interface {
	GetSession(c echo.Context) error
	Logout(c echo.Context) error
	DoNothing(c echo.Context) error
	ServeLoginPage(c echo.Context) error
	ServeErrorPage(c echo.Context) error
	Authorize(c echo.Context) error
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

func (h Handlers) ServeLoginPage(c echo.Context) error {
	accessTokenValue := ""
	if accessToken, err := c.Cookie("access_token"); err == nil {
		accessTokenValue = accessToken.Value
	}
	refreshTokenValue := ""
	if refreshToken, err := c.Cookie("refresh_token"); err == nil {
		refreshTokenValue = refreshToken.Value
	}
	tokensPair, err := h.Service.CheckAndRefreshToken(accessTokenValue, refreshTokenValue, false)
	if tokensPair != nil {
		accessCookie := cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), h.Config)
		refreshCookie := cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), h.Config)
		c.SetCookie(&accessCookie)
		c.SetCookie(&refreshCookie)
	}
	if err == nil {
		return c.Redirect(http.StatusFound, h.Config.App.RedirectAfterSuccess)
	}
	tmpl, err := template.ParseFS(templates, "templates/login.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, apperrors.ErrGeneric.Error())
	}
	data := struct {
		AppName        string
		GitHubEnabled  bool
		DiscordEnabled bool
	}{
		AppName:        h.Config.App.Name,
		GitHubEnabled:  h.Config.Auth.Providers.GitHub.Enabled,
		DiscordEnabled: h.Config.Auth.Providers.Discord.Enabled,
	}
	return tmpl.Execute(c.Response().Writer, data)
}

func (h Handlers) ServeErrorPage(c echo.Context) error {
	tmpl, err := template.ParseFS(templates, "templates/error.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, apperrors.ErrGeneric.Error())
	}
	data := struct {
		AppName string
		Error   string
	}{
		AppName: h.Config.App.Name,
		Error:   c.QueryParam("error"),
	}
	return tmpl.Execute(c.Response().Writer, data)
}

func (h Handlers) Authorize(c echo.Context) error {
	type Body struct {
		AccessToken string   `json:"access_token"`
		Roles       []string `json:"authorized_roles"`
	}
	body := Body{}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	err := h.Service.Authorize(body.AccessToken, body.Roles)
	if err != nil {
		if errors.Is(err, apperrors.ErrAccessTokenExpired) {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error":      apperrors.ErrAccessTokenExpired.Error(),
				"authorized": false,
			})
		}
		if errors.Is(err, apperrors.ErrAccessTokenInvalid) {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error":      apperrors.ErrAccessTokenInvalid.Error(),
				"authorized": false,
			})
		}
		if errors.Is(err, apperrors.ErrRefreshTokenInvalid) {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error":      apperrors.ErrRefreshTokenInvalid.Error(),
				"authorized": false,
			})
		}
		if errors.Is(err, apperrors.ErrNoRoles) {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error":      apperrors.ErrNoRoles.Error(),
				"authorized": false,
			})
		}
		if errors.Is(err, apperrors.ErrUnauthorizedRole) {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error":      apperrors.ErrUnauthorizedRole.Error(),
				"authorized": false,
			})
		}
		if errors.Is(err, apperrors.ErrAccessTokenInvalid) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenInvalid.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	return c.JSON(http.StatusOK, map[string]bool{"authorized": true})
}
