package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxClient interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Exec(ctx context.Context, client pgxClient, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return client.Exec(ctx, sql, arguments...)
}

func QueryRow(ctx context.Context, client pgxClient, sql string, args ...any) pgx.Row {
	return client.QueryRow(ctx, sql, args...)
}
