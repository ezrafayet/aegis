package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"aegix/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() error {
	conf, err := config.ReadConfig("config.json")
	if err != nil {
		return err
	}

	fmt.Println(`
 █████╗ ██╗   ██╗████████╗██╗  ██╗     █████╗ ███████╗ ██████╗ ██╗██╗  ██╗
██╔══██╗██║   ██║╚══██╔══╝██║  ██║    ██╔══██╗██╔════╝██╔════╝ ██║╚██╗██╔╝
███████║██║   ██║   ██║   ███████║    ███████║█████╗  ██║  ███╗██║ ╚███╔╝ 
██╔══██║██║   ██║   ██║   ██╔══██║    ██╔══██║██╔══╝  ██║   ██║██║ ██╔██╗ 
██║  ██║╚██████╔╝   ██║   ██║  ██║    ██║  ██║███████╗╚██████╔╝██║██╔╝ ██╗
╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝    ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝
Drop-in auth service - no SaaS, no lock-in
v0.1.0
	`)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: conf.Auth.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// /!\ issues an access token
	e.POST("/auth/github/callback", func(c echo.Context) error {
		type GithubCallbackRequest struct {
			Code string `json:"code"`
			State string `json:"state"`
		}
		var args GithubCallbackRequest
		if err := c.Bind(&args); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		type GitHubTokenResponse struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
			Scope       string `json:"scope"`
		}
		
		data := map[string]string{
			"client_id":     conf.Auth.Providers.GitHub.ClientID,
			"client_secret": conf.Auth.Providers.GitHub.ClientSecret,
			"code":          args.Code,
			// "redirect_uri":  "http://localhost:3000/auth/callback", // needed ?
			"state":         args.State,
		}
		body, _ := json.Marshal(data)
		
		req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(body))
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get access token"})
		}
		defer resp.Body.Close()

		var tokenResponse GitHubTokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode access token"})
		}

		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/me", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/refresh", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/logout", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/internal/authorize", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/internal/authorize-api-key", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":5666"); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
