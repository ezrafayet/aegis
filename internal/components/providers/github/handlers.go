package github

import (
	"aegix/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OAuthGithubHandlers struct {
	Config  domain.Config
	Service domain.OAuthProviderService
}

var _ domain.OAuthProviderHandlers = OAuthGithubHandlers{}

func NewOAuthGithubHandlers(c domain.Config, s domain.OAuthProviderService) OAuthGithubHandlers {
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
	accessCookie, refreshCookie, err := h.Service.ExchangeCode(body.Code, body.State)
	if err != nil {
		if err.Error() == domain.ErrEarlyAdoptersOnly.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrEarlyAdoptersOnly.Error()})
		}
		if err.Error() == domain.ErrUserBlocked.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrUserBlocked.Error()})
		}
		if err.Error() == domain.ErrUserDeleted.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrUserDeleted.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": domain.ErrGeneric.Error()})
	}
	c.SetCookie(&accessCookie)
	c.SetCookie(&refreshCookie)
	return c.NoContent(http.StatusOK)
}
