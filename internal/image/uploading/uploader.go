package uploading

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/oxtoacart/bpool"

	"github.com/Eviljeks/blurer/internal/hasher"
	"github.com/Eviljeks/blurer/internal/image"
	"github.com/Eviljeks/blurer/internal/image/adding"
	"github.com/Eviljeks/blurer/internal/image/encoding"
	"github.com/Eviljeks/blurer/internal/store"
	"github.com/Eviljeks/blurer/pkg/clock"
)

type Uploader struct {
	adder   *adding.Adder
	buffers *bpool.BufferPool
	store   *store.Store
	hasher  hasher.Hasher
}

func NewUploader(adder *adding.Adder, buffers *bpool.BufferPool, s *store.Store, hasher hasher.Hasher) *Uploader {
	return &Uploader{adder: adder, buffers: buffers, store: s, hasher: hasher}
}

func (u *Uploader) Upload(ctx context.Context, r io.Reader) (*image.Image, error, bool) {
	buf := u.buffers.Get()
	defer u.buffers.Put(buf)

	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err, false
	}
	blob := buf.Bytes()

	filepath := fmt.Sprintf("%s.%s", u.generateFilename(blob), string(encoding.JpegExt))

	foundImg, err := u.store.GetImage(ctx, "", filepath)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err, false
		}
	}
	if foundImg != nil {
		return foundImg, nil, true
	}

	path, err := u.adder.Add(buf, filepath)
	if err != nil {
		return nil, fmt.Errorf("error adding image: %s", err), false
	}

	img := image.Image{
		UUID: uuid.NewString(),
		Path: path,
		TS:   clock.GetCurrentTS(),
	}
	_, err = u.store.SaveImage(ctx, img)
	if err != nil {
		return nil, fmt.Errorf("error saving image: %s", err), false
	}

	return &img, nil, false
}

func (u *Uploader) generateFilename(buf []byte) string {
	filename := u.hasher.Hash(buf)

	return filename
}
