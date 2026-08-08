package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/middleware"
	linkedin "github.com/prashanth30-n/Pubo/internal/model/Linkedin"
	"github.com/prashanth30-n/Pubo/internal/server"
	"github.com/prashanth30-n/Pubo/internal/service"
)

type LinkedinHandler struct {
	Handler
	LinkedinService *service.LinkedinService
}

func NewLinkedinHandler(s *server.Server, LinkedinService *service.LinkedinService) *LinkedinHandler {
	return &LinkedinHandler{
		Handler:         NewHandler(s),
		LinkedinService: LinkedinService,
	}
}
func (h *LinkedinHandler) ConnectLinkedin(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *linkedin.CreateLinkedinConnectionPayload) (*linkedin.LinkedinConnectedAccount, error) {
			userId := middleware.GetUserID(c)
			return h.LinkedinService.CreateLinkedinConnection(c, userId, payload)
		},
		http.StatusCreated,
		&linkedin.CreateLinkedinConnectionPayload{},
	)(c)

}
