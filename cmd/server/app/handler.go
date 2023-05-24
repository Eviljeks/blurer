package app

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/oxtoacart/bpool"

	"github.com/Eviljeks/blurer/internal/api/route"
	"github.com/Eviljeks/blurer/internal/hasher"
	"github.com/Eviljeks/blurer/internal/image/adding"
	"github.com/Eviljeks/blurer/internal/image/imageprocessing"
	"github.com/Eviljeks/blurer/internal/image/uploading"
	"github.com/Eviljeks/blurer/internal/store"
)

type Storage interface {
	Open(src string) (io.ReadWriteCloser, error)
	Write(r io.Reader, filepath string) error
}

func NewHandler(cfg *Config, storage Storage, conn *pgx.Conn) *gin.Engine {
	s := store.NewStore(conn)
	blurer := imageprocessing.NewBlurer()
	adder := adding.NewAdder(storage)
	uploader := uploading.NewUploader(adder, bpool.NewBufferPool(cfg.bufferPoolSize), s, hasher.NewSha1Hasher())

	r := gin.Default()
	route.UploadImage(r, uploader)
	route.ListImages(r, s)

	blurHandler := route.NewBlurHandler(blurer, s, storage, uploader)

	blurHandler.BlurUploadedImage(r)
	blurHandler.BlurImage(r)

	route.ListImagesBlurred(r, s)

	r.Static("/resources/image", cfg.dstImagePath)

	return r
}
