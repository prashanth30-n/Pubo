package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/handler"
	"github.com/prashanth30-n/Pubo/internal/middleware"
)

func registerSignUpRoute(r *echo.Group, h *handler.SignUpHandler, auth *middleware.AuthMiddleware) {
	signup := r.Group("/signup")
	signup.Use(auth.RequireAuth)
	signup.POST("", h.SignUp)
}
