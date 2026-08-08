package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/handler"
	"github.com/prashanth30-n/Pubo/internal/middleware"
)

func registerBlueskyRoutes(r *echo.Group, h *handler.BlueskyHandler, auth *middleware.AuthMiddleware) {
	bluesky := r.Group("/bluesky")
	bluesky.Use(auth.RequireAuth)
	bluesky.POST("/connect", h.ConnectBluesky)
}
