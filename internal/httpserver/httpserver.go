package httpserver

import (
	"fmt"
	"net/http"

	"aegix/internal/registry"
	"aegix/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() error {
	c, err := config.ReadConfig("config.json")
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
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:           3600,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     c.Auth.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	r := registry.NewRegistry(c)

	r.GitHubRouter.AttachRoutes(e)

	e.GET("/me", func(c echo.Context) error {
		// decode and return the jwt, refresh if needed
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/refresh", func(c echo.Context) error {
		// refresh the jwt
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/logout", func(c echo.Context) error {
		// delete the jwt
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/authorize", func(c echo.Context) error {
		// ask auth service if a jwt is valid, and get user's details from jwt
		return c.String(http.StatusOK, "Hello, World!")
	})

	// must also retrieve and set
	// e.POST("/authorize-api-token", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	// get/set user metadata

	if err := e.Start(":5666"); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
