package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/middleware"
	"github.com/prashanth30-n/Pubo/internal/model/media"
	"github.com/prashanth30-n/Pubo/internal/server"
	"github.com/prashanth30-n/Pubo/internal/service"
)

type MediaHandler struct {
	Handler
	service service.MediaService
}

func NewMediaHandler(s *server.Server, MediaService *service.MediaService) *MediaHandler {
	return &MediaHandler{
		Handler: NewHandler(s),
		service: *MediaService,
	}
}
func (h *MediaHandler) Upload(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *media.UploadMediaPayload) (*media.Asset, error) {
			userID := middleware.GetUserID(c)
			fileheader, err := c.FormFile("file")
			if err != nil || fileheader == nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "file is required")
			}
			return h.service.Upload(c, userID, fileheader, payload.Folder)
		},
		http.StatusCreated,
		&media.UploadMediaPayload{},
	)(c)
}
func (h *MediaHandler) Get(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *EmptyRequest) (*media.Asset, error) {
			id := c.Param("id")
			return h.service.Get(c, id)
		},
		http.StatusOK,
		&EmptyRequest{},
	)(c)
}
func (h *MediaHandler) List(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *EmptyRequest) ([]media.Asset, error) {
			userID := middleware.GetUserID(c)
			limit := 20
			offset := 0
			if v := c.QueryParam("limit"); v != "" {
				if n, err := strconv.Atoi(v); err == nil && n > 0 {
					limit = n
				}
			}
			if v := c.QueryParam("offset"); v != "" {
				if n, err := strconv.Atoi(v); err == nil && n >= 0 {
					offset = n
				}
			}
			return h.service.List(c, userID, limit, offset)
		},
		http.StatusOK,
		&EmptyRequest{},
	)(c)
}

// func(h*MediaHandler) Delete(c echo.Context) error{
// 	return Handle(
// 		h.Handler,
// 		func(c echo.Context,payload *EmptyRequest)error{
// 			id:=c.Param("id")
// 			return h.service.Delete(c,id)
// 		},
// 		http.StatusNoContent,
// 		&EmptyRequest{},
// 	)(c)
// }
