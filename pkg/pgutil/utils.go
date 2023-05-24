package pgutil

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/Eviljeks/blurer/pkg"
)

func SB() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func Connect(ctx context.Context, connString string) (*pgx.Conn, error) {
	var conn *pgx.Conn
	err := pkg.NewWaiter(time.Second, uint8(5)).Wait(ctx, func(ctx context.Context) error {
		var err error
		conn, err = pgx.Connect(ctx, connString)
		if err != nil {
			if strings.Contains(err.Error(), "dial tcp") {
				return pkg.ErrNotReadyYet
			}

			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
