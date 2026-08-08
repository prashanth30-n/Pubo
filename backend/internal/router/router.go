package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/prashanth30-n/Pubo/internal/handler"
	"github.com/prashanth30-n/Pubo/internal/middleware"
	v1 "github.com/prashanth30-n/Pubo/internal/router/v1"
	"github.com/prashanth30-n/Pubo/internal/server"
	"github.com/prashanth30-n/Pubo/internal/service"
	"golang.org/x/time/rate"
)

func NewRouter(s *server.Server, h *handler.Handlers, services *service.Services) *echo.Echo {
	middlewares := middleware.NewMiddlewares(s)
	router := echo.New()
	router.HTTPErrorHandler = middlewares.Global.GlobalErrorHanlder
	router.Use(
		echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
			Store: echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20)),
			DenyHandler: func(c echo.Context, identifier string, err error) error {
				//record rate limit middlewares
				if rateLimitMiddleware := middlewares.RateLimit; rateLimitMiddleware != nil {
					rateLimitMiddleware.RecordRateLimitHit(c.Path())
				}
				s.Logger.Warn().Str("request_id", middleware.GetRequestID(c)).
					Str("identifier", identifier).Str("path", c.Path()).Str("method", c.Request().Method).Str("ip", c.RealIP()).Msg("rate limit exceeded")
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			},
		}),
		middlewares.Global.CORS(),
		middlewares.Global.Secure(),
		middleware.RequestID(),
		middlewares.Tracing.NewRelicMiddleware(),
		middlewares.Tracing.EnhanceTracing(),
		middlewares.ContextEnhancer.EnhanceContext(),
		middlewares.Global.RequestLogger(),
		middlewares.Global.Recover(),
	)
	registerSystemRoutes(router, h)
	apiV1 := router.Group("/api/v1")
	v1.RegisterV1Routes(apiV1, h, middlewares)
	return router
}
