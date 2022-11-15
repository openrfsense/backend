package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Runs a query which does not return a result.
func Do(ctx context.Context, sql string, args ...any) error {
	_, err := pg.Exec(ctx, sql, args...)
	return err
}

// Runs a query which returns a single row and converts it to T using
// pgx.RowToAddrOfStructByName. If no records were returned by the query,
// a nil pointer to T is returned, along with the pgx error.
func Single[T any](ctx context.Context, sql string, args ...any) (*T, error) {
	rows, err := pg.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectOneRow(rows, RowToAddrOfStructByName[T])
}

// Runs a query which returns multiple rows and collects them in a slice of type []T
// using pgx.RowToStructByName. If no records were returned by the query,
// an empty slice is returned, along with the pgx error.
func Multiple[T any](ctx context.Context, sql string, args ...any) ([]T, error) {
	empty := []T{}
	rows, err := pg.Query(ctx, sql, args...)
	if err != nil {
		return empty, err
	}

	return pgx.CollectRows(rows, RowToStructByName[T])
}
