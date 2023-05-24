package store

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Eviljeks/blurer/internal/image"
	"github.com/Eviljeks/blurer/pkg/pgutil"
)

const (
	BlurImageTable           = "image_blurred"
	BlurImageColumnUUID      = "uuid"
	BlurImageColumnImageUUID = "image_uuid"
	BlurImageColumnX0        = "x_0"
	BlurImageColumnY0        = "y_0"
	BlurImageColumnX1        = "x_1"
	BlurImageColumnY1        = "y_1"
	BlurImageColumnTS        = "ts"
)

func (s *Store) SaveBlurImage(ctx context.Context, bImage image.BlurImage) (bool, error) {
	sb := pgutil.SB()
	sql, args, err := sb.Insert(BlurImageTable).
		Columns(
			BlurImageColumnUUID,
			BlurImageColumnImageUUID,
			BlurImageColumnX0,
			BlurImageColumnY0,
			BlurImageColumnX1,
			BlurImageColumnY1,
			BlurImageColumnTS,
		).
		Values(
			bImage.UUID,
			bImage.ImageUUID,
			bImage.X0,
			bImage.Y0,
			bImage.X1,
			bImage.Y1,
			bImage.TS,
		).
		ToSql()

	if err != nil {
		return false, errors.Wrapf(err, "blurImage[%s]", bImage.ImageUUID)
	}

	tag, err := s.Conn.Exec(ctx, sql, args...)
	if err != nil {
		return false, errors.Wrapf(err, "blurImage[%s]", bImage.ImageUUID)
	}

	return tag.RowsAffected() > 0, nil
}

func (s *Store) ListBlurImages(ctx context.Context, limit uint64, imgUUID string) ([]*image.BlurImage, error) {
	sb := pgutil.SB().
		Select(
			BlurImageColumnUUID,
			BlurImageColumnImageUUID,
			BlurImageColumnX0,
			BlurImageColumnY0,
			BlurImageColumnX1,
			BlurImageColumnY1,
			BlurImageColumnTS,
		).
		From(BlurImageTable)

	if imgUUID != "" {
		sb = sb.Where("image_uuid = ?", imgUUID)
	}

	sql, args, err := sb.
		OrderBy(BlurImageColumnTS + " DESC").
		Limit(limit).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := s.Conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	results := make([]*image.BlurImage, 0, limit)

	for rows.Next() {
		var bImg image.BlurImage

		if err = rows.Scan(
			&bImg.UUID,
			&bImg.ImageUUID,
			&bImg.X0,
			&bImg.Y0,
			&bImg.X1,
			&bImg.Y1,
			&bImg.TS,
		); err != nil {
			return nil, err
		}

		results = append(results, &bImg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
