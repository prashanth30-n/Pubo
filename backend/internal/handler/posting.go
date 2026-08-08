package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/middleware"
	"github.com/prashanth30-n/Pubo/internal/server"
	"github.com/prashanth30-n/Pubo/internal/service"
)

type PostHandler struct {
	Handler
	service *service.PostService
}

func NewPostHandler(s *server.Server, postService *service.PostService) *PostHandler {
	return &PostHandler{Handler: NewHandler(s), service: postService}
}
func (h *PostHandler) Publish(c echo.Context) error {
	var request service.PublishRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	response, err := h.service.Publish(c.Request().Context(), middleware.GetUserID(c), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}
func (h *PostHandler) SaveDraft(c echo.Context) error {
	var request service.DraftRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	post, err := h.service.SaveDraft(c.Request().Context(), middleware.GetUserID(c), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, post)
}
func (h *PostHandler) List(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	posts, err := h.service.List(c.Request().Context(), middleware.GetUserID(c), c.QueryParam("status"), limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, posts)
}
