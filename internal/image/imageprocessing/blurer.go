package imageprocessing

import (
	"image"
	"image/draw"

	"github.com/esimov/stackblur-go"
)

type Blurer struct{}

func NewBlurer() *Blurer {
	return &Blurer{}
}

func (b *Blurer) Blur(srcImage image.Image, dstRect BlurRectangle, blurRad uint32) (image.Image, error) {
	bounds := srcImage.Bounds()
	m := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(m, m.Bounds(), srcImage, bounds.Min, draw.Src)

	res, err := stackblur.Process(srcImage, blurRad)
	if err != nil {
		return nil, err
	}

	draw.Draw(m, dstRect.Rectangle(), res, dstRect.GetMin(), draw.Src)

	return m, nil
}

type BlurRectangle struct {
	min, max image.Point
}

func (br *BlurRectangle) Rectangle() image.Rectangle {
	return image.Rect(br.min.X, br.min.Y, br.max.X, br.max.Y)
}

func (br *BlurRectangle) GetMin() image.Point {
	return br.min
}

func NewBlurRectangle(x0, y0, x1, y1 int, bounds image.Rectangle) BlurRectangle {
	if x0 < bounds.Min.X || x0 > bounds.Max.X {
		x0 = 0
	}

	if y0 < bounds.Min.Y || y0 > bounds.Max.Y {
		y0 = 0
	}

	if x1 == 0 {
		x1 = bounds.Max.X
	} else if x1 < 0 {
		x1 = bounds.Max.X + x1
	}
	if y1 == 0 {
		y1 = bounds.Max.Y
	} else if y1 < 0 {
		y1 = bounds.Max.Y + y1
	}

	if x1 < bounds.Min.X || x1 > bounds.Max.X {
		x1 = bounds.Max.X
	}

	if y1 < bounds.Min.Y || y1 > bounds.Max.Y {
		y1 = bounds.Max.Y
	}

	return BlurRectangle{
		min: image.Point{X: x0, Y: y0},
		max: image.Point{X: x1, Y: y1},
	}
}
