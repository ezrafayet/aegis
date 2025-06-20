package github

import (
	"net/http"
	"othnx/internal/components/providers/providersports"
	"othnx/internal/domain"

	"github.com/labstack/echo/v4"
)

type OAuthGithubHandlers struct {
	Config  domain.Config
	Service providersports.OAuthProviderService
}

var _ providersports.OAuthProviderHandlers = OAuthGithubHandlers{}

func NewOAuthGithubHandlers(c domain.Config, s providersports.OAuthProviderService) OAuthGithubHandlers {
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
	if accessCookie != nil && refreshCookie != nil {
		c.SetCookie(accessCookie)
		c.SetCookie(refreshCookie)
	}
	return c.NoContent(http.StatusOK)
}
