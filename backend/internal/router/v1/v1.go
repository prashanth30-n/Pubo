package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/handler"
	"github.com/prashanth30-n/Pubo/internal/middleware"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers, middlewares *middleware.MiddleWares) {
	registerQuoteRoutes(router, handlers.Quotes, middlewares.Auth)
	registerMediaRoutes(router, handlers.Media, middlewares.Auth)
	registerBlueskyRoutes(router, handlers.Bluesky, middlewares.Auth)
	registerSignUpRoute(router, handlers.Signup, middlewares.Auth)
	RegisterLinkedinRoutes(router, handlers.Linkedin, middlewares.Auth)
	registerPostRoutes(router, handlers.Posts, middlewares.Auth)
}
