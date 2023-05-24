package route

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Eviljeks/blurer/internal/api"
	"github.com/Eviljeks/blurer/internal/image/uploading"
	"github.com/Eviljeks/blurer/internal/store"
)

type UploadImageParams struct {
	File *multipart.FileHeader `form:"file"`
}

type ListImagesParams struct {
	Limit uint64 `form:"limit,default=20"`
}

func UploadImage(r gin.IRouter, u *uploading.Uploader) {
	r.POST("/image", func(ctx *gin.Context) {
		var params UploadImageParams

		if err := ctx.ShouldBind(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		if ext := filepath.Ext(params.File.Filename); ext != ".jpg" && ext != ".jpeg" {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(fmt.Errorf("ext is not jpg: %s", ext)))
			return
		}

		file, err := params.File.Open()
		defer file.Close()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		img, err, _ := u.Upload(ctx, file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(fmt.Errorf("error uploading image: %s", err)))
			return
		}

		ctx.JSON(http.StatusOK, api.OK(img))

		return
	})
}

func ListImages(r gin.IRouter, s *store.Store) {
	r.GET("/image", func(ctx *gin.Context) {
		var params ListImagesParams

		if err := ctx.ShouldBind(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		images, err := s.ListImages(ctx, params.Limit)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(fmt.Errorf("error listing images: %s", err)))
			return
		}

		ctx.JSON(http.StatusOK, api.OK(images))

		return
	})
}
