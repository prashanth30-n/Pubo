package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/middleware"
	bluesky "github.com/prashanth30-n/Pubo/internal/model/Bluesky"
	"github.com/prashanth30-n/Pubo/internal/server"
	"github.com/prashanth30-n/Pubo/internal/service"
)

type BlueskyHandler struct {
	Handler
	blueskyService *service.BlueskyService
}

func NewBlueskyHandler(s *server.Server, blueskyService *service.BlueskyService) *BlueskyHandler {
	return &BlueskyHandler{
		Handler:        NewHandler(s),
		blueskyService: blueskyService,
	}
}
func (h *BlueskyHandler) ConnectBluesky(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *bluesky.CreateBlueskyConnectionPayload) (*bluesky.ConnectedAccounts, error) {
			userID := middleware.GetUserID(c)
			return h.blueskyService.ConnectBluesky(c, userID, payload)
		},
		http.StatusCreated,
		&bluesky.CreateBlueskyConnectionPayload{},
	)(c)
}
