package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type OpenAPIHandler struct {
	Handler
}

func NewOpenAPIHanlder(s *server.Server) *OpenAPIHandler {
	return &OpenAPIHandler{
		Handler: NewHandler(s),
	}
}

func (h *OpenAPIHandler) ServerOpenAPIUI(c echo.Context) error {
	templateBytes, err := os.ReadFile("static/openapi.html")
	c.Response().Header().Set("Cache-Control", "no-cache")
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI UI template: %w", err)
	}
	templateString := string(templateBytes)
	err = c.HTML(http.StatusOK, templateString)
	if err != nil {
		return fmt.Errorf("failed to write HTML response: %w", err)
	}
	return nil
}
