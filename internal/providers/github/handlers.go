package github

import (
	"aegix/internal/domain"
	"aegix/internal/providers"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OAuthGithubHandlers struct {
	Config  domain.Config
	Service providers.OAuthProviderService
}

var _ providers.OAuthProviderHandlers = OAuthGithubHandlers{}

func NewOAuthGithubHandlers(c domain.Config, s providers.OAuthProviderService) OAuthGithubHandlers {
	return OAuthGithubHandlers{
		Config: c,
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
	// get code and state fromn the json body NOT FORM
	type ExchangeCodeRequest struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
	var body ExchangeCodeRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	accessToken, err := h.Service.ExchangeCode(body.Code, body.State)
	if err != nil {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an error occurred"})
	}
	// todo: save access token to cookie
	return c.JSON(http.StatusOK, map[string]string{"access_token": accessToken})
}
