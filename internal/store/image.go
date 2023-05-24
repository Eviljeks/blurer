package store

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Eviljeks/blurer/internal/image"
	"github.com/Eviljeks/blurer/pkg/pgutil"
)

const (
	ImageTable      = "image"
	ImageColumnUUID = "uuid"
	ImageColumnPath = "path"
	ImageColumnTS   = "ts"
)

func (s *Store) SaveImage(ctx context.Context, image image.Image) (bool, error) {
	sb := pgutil.SB()
	sql, args, err := sb.Insert(ImageTable).
		Columns(
			ImageColumnUUID,
			ImageColumnPath,
			ImageColumnTS,
		).
		Values(
			image.UUID,
			image.Path,
			image.TS,
		).
		ToSql()

	if err != nil {
		return false, errors.Wrapf(err, "image[%s]", image.UUID)
	}

	tag, err := s.Conn.Exec(ctx, sql, args...)
	if err != nil {
		return false, errors.Wrapf(err, "image[%s]", image.UUID)
	}

	return tag.RowsAffected() > 0, nil
}

func (s *Store) GetImage(ctx context.Context, uuid string, path string) (*image.Image, error) {
	sb := pgutil.SB().
		Select(
			ImageColumnUUID,
			ImageColumnPath,
			ImageColumnTS,
		).
		From(ImageTable)

	if path == "" && uuid == "" {
		return nil, errors.New("'uuid' or 'path' should be provided")
	}

	if path != "" {
		sb = sb.Where("path = ?", path)
	}

	if uuid != "" {
		sb = sb.Where("uuid = ?", uuid)
	}

	sql, args, err := sb.ToSql()

	if err != nil {
		return nil, err
	}

	img := image.Image{}

	err = s.Conn.QueryRow(ctx, sql, args...).Scan(
		&img.UUID,
		&img.Path,
		&img.TS,
	)

	if err != nil {
		return nil, err
	}

	return &img, nil
}

func (s *Store) ListImages(ctx context.Context, limit uint64) ([]*image.Image, error) {
	sql, args, err := pgutil.SB().
		Select(
			ImageColumnUUID,
			ImageColumnPath,
			ImageColumnTS,
		).
		From(ImageTable).
		OrderBy(ImageColumnTS + " DESC").
		Limit(limit).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := s.Conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	results := make([]*image.Image, 0, limit)

	for rows.Next() {
		var img image.Image

		if err = rows.Scan(
			&img.UUID,
			&img.Path,
			&img.TS,
		); err != nil {
			return nil, err
		}

		results = append(results, &img)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
