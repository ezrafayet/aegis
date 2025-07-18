package httpserver

import (
	"fmt"
	"net/http"

	"aegis/internal/infrastructure/config"
	"aegis/internal/infrastructure/database"
	"aegis/internal/registry"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Version = "dev"

func Start() error {
	fmt.Println(`
 █████╗ ███████╗ ██████╗ ██╗███████╗
██╔══██╗██╔════╝██╔════╝ ██║██╔════╝
███████║█████╗  ██║  ███╗██║███████╗
██╔══██║██╔══╝  ██║   ██║██║╚════██║
██║  ██║███████╗╚██████╔╝██║███████║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚══════╝
Drop-in auth service - no SaaS, no lock-in
` + Version + `
	`)
	c, err := config.Read("config.json")
	if err != nil {
		return err
	}

	db, err := database.Connect(c)
	if err != nil {
		return err
	}

	if err := database.Migrate(db); err != nil {
		return err
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		HSTSMaxAge:         3600,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     c.App.CorsAllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	r, err := registry.NewRegistry(c, db)
	if err != nil {
		return err
	}

	group := e.Group("/auth")

	group.GET("/me", r.Handlers.GetSession, r.Middlewares.CheckAndRefreshToken)
	group.GET("/refresh", r.Handlers.DoNothing, r.Middlewares.CheckAndForceRefreshToken)
	group.GET("/logout", r.Handlers.Logout)
	group.GET("/health", r.Handlers.DoNothing)
	group.POST("/authorize-access-token", r.Handlers.Authorize)

	if c.LoginPage.Enabled {
		e.GET(c.LoginPage.FullPath, r.Handlers.ServeLoginPage)
	}

	if c.ErrorPage.Enabled {
		e.GET(c.ErrorPage.FullPath, r.Handlers.ServeErrorPage)
	}

	for _, provider := range r.Providers {
		group.GET(fmt.Sprintf("/%s", provider.Name), provider.Handlers.GetAuthURL, provider.Middlewares.CheckAuthEnabled)
		group.GET(fmt.Sprintf("/%s/callback", provider.Name), provider.Handlers.ExchangeCode, provider.Middlewares.CheckAuthEnabled)
	}

	return e.Start(fmt.Sprintf(":%d", c.App.Port))
}
