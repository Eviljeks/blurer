package route

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Eviljeks/blurer/internal/api"
	img "github.com/Eviljeks/blurer/internal/image"
	"github.com/Eviljeks/blurer/internal/image/encoding"
	"github.com/Eviljeks/blurer/internal/image/imageprocessing"
	"github.com/Eviljeks/blurer/internal/image/uploading"
	"github.com/Eviljeks/blurer/internal/storage"
	"github.com/Eviljeks/blurer/internal/store"
)

type BlurHandler struct {
	blurer     *imageprocessing.Blurer
	s          *store.Store
	imgStorage *storage.Storage
	uploader   *uploading.Uploader
}

func NewBlurHandler(blurer *imageprocessing.Blurer, s *store.Store, imgStorage *storage.Storage, uploader *uploading.Uploader) *BlurHandler {
	return &BlurHandler{blurer: blurer, s: s, imgStorage: imgStorage, uploader: uploader}
}

type BlurUploadedImageParams struct {
	X0     int `form:"x_0" json:"x_0"`
	Y0     int `form:"y_0" json:"y_0"`
	X1     int `form:"x_1" json:"x_1"`
	Y1     int `form:"y_1" json:"y_1"`
	Radius int `form:"radius" json:"radius"`
}

func (bh *BlurHandler) BlurUploadedImage(r gin.IRouter) {
	r.PUT("/image/:uuid/blur", func(ctx *gin.Context) {
		imgUUID := ctx.Param("uuid")
		var params BlurUploadedImageParams

		if err := ctx.ShouldBind(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		i, err := bh.s.GetImage(ctx, imgUUID, "")
		if err != nil {
			ctx.JSON(http.StatusNotFound, api.ErrorNotFound(fmt.Errorf("image not found: %s, %s", imgUUID, err)))
			return
		}

		blurImage, err := bh.doBlur(ctx, i, params.Radius, params.X0, params.Y0, params.X1, params.Y1)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		ctx.JSON(http.StatusOK, api.OK(blurImage))

		return
	})
}

type BlurImageParams struct {
	File   *multipart.FileHeader `form:"file"`
	X0     int                   `form:"x_0"`
	Y0     int                   `form:"y_0"`
	X1     int                   `form:"x_1"`
	Y1     int                   `form:"y_1"`
	Radius int                   `form:"radius"`
}

func (bh *BlurHandler) BlurImage(r gin.IRouter) {
	r.PUT("/image/blur", func(ctx *gin.Context) {
		var params BlurImageParams

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

		i, err, _ := bh.uploader.Upload(ctx, file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(fmt.Errorf("error uploading image: %s", err)))
			return
		}

		blurImage, err := bh.doBlur(ctx, i, params.Radius, params.X0, params.Y0, params.X1, params.Y1)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		ctx.JSON(http.StatusOK, api.OK(blurImage))

		return
	})
}

func (bh *BlurHandler) doBlur(
	ctx context.Context,
	i *img.Image,
	r int,
	x0 int,
	y0 int,
	x1 int,
	y1 int,
) (*img.BlurImage, error) {
	imgFile, err := bh.imgStorage.Open(i.Path)
	if err != nil {
		return nil, fmt.Errorf("could not open image file: %s, %s", imgFile, err)
	}

	srcImage, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("could not decode image file: %s, %s", imgFile, err)
	}

	blurRectangle := imageprocessing.NewBlurRectangle(x0, y0, x1, y1, srcImage.Bounds())

	bluredImage, err := bh.blurer.Blur(
		srcImage,
		blurRectangle,
		uint32(r),
	)
	if err != nil {
		return nil, fmt.Errorf("could not blur image: %s, %s", imgFile, err)
	}

	buf := bytes.NewBuffer(nil)

	encoder := encoding.NewJpegEncoder()
	err = encoder.Encode(buf, bluredImage)
	if err != nil {
		return nil, fmt.Errorf("could not encode image: %s, %s", imgFile, err)
	}

	newImg, err, found := bh.uploader.Upload(ctx, buf)
	if err != nil {
		return nil, fmt.Errorf("could not upload image: %s, %s", i.Path, err)
	}

	blurImage := img.BlurImage{
		UUID:      newImg.UUID,
		ImageUUID: i.UUID,
		X0:        blurRectangle.Rectangle().Min.X,
		Y0:        blurRectangle.Rectangle().Min.Y,
		X1:        blurRectangle.Rectangle().Max.X,
		Y1:        blurRectangle.Rectangle().Max.Y,
		TS:        newImg.TS,
	}

	if !found {
		_, err = bh.s.SaveBlurImage(ctx, blurImage)
		if err != nil {
			return nil, fmt.Errorf("error saving blurImage: %s, %s", blurImage.UUID, err)
		}
	}

	return &blurImage, nil
}
