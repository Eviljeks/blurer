package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"

	"github.com/Eviljeks/blurer/internal/api/route"
	"github.com/Eviljeks/blurer/internal/hasher"
	"github.com/Eviljeks/blurer/internal/image/adding"
	"github.com/Eviljeks/blurer/internal/image/imageprocessing"
	"github.com/Eviljeks/blurer/internal/image/uploading"
	"github.com/Eviljeks/blurer/internal/storage"
	"github.com/Eviljeks/blurer/internal/store"
	"github.com/Eviljeks/blurer/pkg/pgutil"
)

type Config struct {
	dstImagePrefix string
	bufferPoolSize int
}

func defaultConfig() *Config {
	return &Config{
		dstImagePrefix: "./storage/image",
		bufferPoolSize: 10,
	}
}

func main() {
	ctx := context.Background()
	conn, err := pgutil.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.Fatalf("db connection faild, err: %s", err.Error())
	}
	logrus.Infoln("db connected")
	defer conn.Close(ctx)

	conf := defaultConfig()

	s := store.NewStore(conn)

	blurer := imageprocessing.NewBlurer()
	imgStorage := storage.NewStorage(conf.dstImagePrefix)
	adder := adding.NewAdder(imgStorage)
	uploader := uploading.NewUploader(adder, bpool.NewBufferPool(conf.bufferPoolSize), s, hasher.NewSha1Hasher())

	r := gin.Default()
	route.UploadImage(r, uploader)
	route.ListImages(r, s)

	blurHandler := route.NewBlurHandler(blurer, s, imgStorage, uploader)

	blurHandler.BlurUploadedImage(r)
	blurHandler.BlurImage(r)

	route.ListImagesBlurred(r, s)

	r.Static("/resources/image", conf.dstImagePrefix)

	go func() {
		sErr := r.Run(":3000")
		if sErr != nil {
			logrus.Fatalf("failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logrus.Print("Server received shutdown signal")
}
