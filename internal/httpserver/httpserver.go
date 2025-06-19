package httpserver

import (
	"fmt"
	"net/http"

	"othnx/internal/domain"
	"othnx/internal/registry"
	"othnx/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Start() error {
	fmt.Println(`
 ██████╗ ████████╗██╗  ██╗███╗   ██╗██╗  ██╗     █████╗ ██╗   ██╗████████╗██╗  ██╗
██╔═══██╗╚══██╔══╝██║  ██║████╗  ██║╚██╗██╔╝    ██╔══██╗██║   ██║╚══██╔══╝██║  ██║
██║   ██║   ██║   ███████║██╔██╗ ██║ ╚███╔╝     ███████║██║   ██║   ██║   ███████║
██║   ██║   ██║   ██╔══██║██║╚██╗██║ ██╔██╗     ██╔══██║██║   ██║   ██║   ██╔══██║
╚██████╔╝   ██║   ██║  ██║██║ ╚████║██╔╝ ██╗    ██║  ██║╚██████╔╝   ██║   ██║  ██║
 ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝
Drop-in auth service - no SaaS, no lock-in
v0.x.x (needs to be injected)
	`)
	c, err := config.ReadConfig("config.json")
	if err != nil {
		return err
	}

	db, err := gorm.Open(postgres.Open(c.DB.PostgresURL))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	} else {
		fmt.Println("Connected to database")
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
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
		AllowOrigins:     c.App.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	r := registry.NewRegistry(c, db)
	r.GitHubRouter.AttachRoutes(e)
	r.AuthRouter.AttachRoutes(e)

	return e.Start(fmt.Sprintf(":%d", c.App.Port))
}
