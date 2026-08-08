package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/handler"
	"github.com/prashanth30-n/Pubo/internal/middleware"
)

func RegisterLinkedinRoutes(r *echo.Group, h *handler.LinkedinHandler, auth *middleware.AuthMiddleware) {
	linkedin := r.Group("/linkedin")
	linkedin.Use(auth.RequireAuth)
	linkedin.POST("/connect", h.ConnectLinkedin)
}
