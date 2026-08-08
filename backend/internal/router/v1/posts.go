package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/handler"
	"github.com/prashanth30-n/Pubo/internal/middleware"
)

func registerPostRoutes(r *echo.Group, h *handler.PostHandler, auth *middleware.AuthMiddleware) {
	posts := r.Group("/posts")
	posts.Use(auth.RequireAuth)
	posts.POST("/drafts", h.SaveDraft)
	posts.GET("", h.List)
	posts.POST("/publish", h.Publish)
}
