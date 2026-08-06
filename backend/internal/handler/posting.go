package handler

import (
		"io"
	"net/http"
	"encoding/json"

	"github.com/PatibandlaVenkat/Pubo/internal/middleware"
	"github.com/labstack/echo/v4"
)
type CreatePostRequest struct{
	Content string `form:"content"`
	AccountIDs []string `form:"account_ids"`
}
func(h *Handler) CreatePost(c echo.Context) error{
	userID:=middleware.GetUserID(c)
	content:=c.FormValue("content")
	var accountIDs []string
	if err:=json.Unmarshal([]byte(c.FormValue("account_ids")),&accountIDs); err!=nil{
		return echo.NewHTTPError(400,"invalid_account_ids")
	}
	form,err:=c.MultipartForm()
	if err!=nil{
		return echo.NewHTTPError(400,"invalid_form_data")
	}
	fileHeaders:=form.File["images"]
	var images []service.ImageFile
	for _,fh:=range fileHeaders{
		src,err:=fh.Open()
		if err!=nil{
			continue
		}
		data,err:=io.ReadAll(src)
		src.Close()
		if err!=nil{
			continue
		}
		images=apppend(images,service.ImageFile{
			Name:fh.Filename,
			Data:data,
			 ContentType:fh.Header.Get("Content-Type"),
		     Size:        fh.Size,
		})
	}
	post,err:=h.postService.CreatePost(c.Request().Context(),userID,content,accountIDs,images)
	if err!=nil{
		return echo.NewHTTPError(500,err.Error())
	}
return c.JSON(http.StatusCreated,post)
}