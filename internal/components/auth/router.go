package auth

import "github.com/labstack/echo/v4"

type AuthRouterInterface interface {
	AttachRoutes(e *echo.Echo)
}

type AuthRouter struct {
	Handlers    AuthHandlersInterface
	Middlewares AuthMiddlewareInterface
}

var _ AuthRouterInterface = &AuthRouter{}

func NewAuthRouter(h AuthHandlersInterface, m AuthMiddlewareInterface) AuthRouter {
	return AuthRouter{
		Handlers:    h,
		Middlewares: m,
	}
}

func (r AuthRouter) AttachRoutes(e *echo.Echo) {
	group := e.Group("/auth", r.Middlewares.CheckAndRefreshToken)
	group.GET("/me", r.Handlers.GetSession)
	group.GET("/refresh", r.Handlers.DoNothing)
	e.GET("/auth/logout", r.Handlers.Logout)
	e.GET("/auth/health", r.Handlers.DoNothing)
}
