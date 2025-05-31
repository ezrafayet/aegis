package httpserver

import (
	"fmt"
	"net/http"

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

	if err := e.Start(":5666"); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
