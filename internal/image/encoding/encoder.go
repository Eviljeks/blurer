package encoding

import (
	"image"
	"image/jpeg"
	"io"
)

type Ext string

const JpegExt Ext = "jpeg"

type JpegEncoder struct{}

func NewJpegEncoder() *JpegEncoder {
	return &JpegEncoder{}
}

func (js *JpegEncoder) Encode(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
}
