package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxClient interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func Exec(ctx context.Context, sql string, client pgxClient, arguments ...any) (pgconn.CommandTag, error) {
	return client.Exec(ctx, sql, arguments)
}
