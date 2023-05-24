package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Eviljeks/blurer/internal/api"
	"github.com/Eviljeks/blurer/internal/store"
)

type ListImagesBlurredParams struct {
	Limit     uint64 `form:"limit,default=20"`
	ImageUUID string `form:"image_uuid"`
}

func ListImagesBlurred(r gin.IRouter, s *store.Store) {
	r.GET("/image/:uuid/blur", func(ctx *gin.Context) {
		imgUUID := ctx.Param("uuid")
		var params ListImagesBlurredParams

		if err := ctx.ShouldBind(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		bImages, err := s.ListBlurImages(ctx, params.Limit, imgUUID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(fmt.Errorf("error listing blur images: %s", err)))
			return
		}

		ctx.JSON(http.StatusOK, api.OK(bImages))

		return
	})
}
