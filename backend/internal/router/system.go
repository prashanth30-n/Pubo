package router

import (
	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/handler"
)

func registerSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/status", h.Health.CheckHealth)
	r.Static("/static", "static")
	r.GET("/docs", h.OpenAPI.ServerOpenAPIUI)
}
