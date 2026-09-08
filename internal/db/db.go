package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}

func NextID(ctx context.Context, pool *pgxpool.Pool, table string) (string, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if err := pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return "", err
	}
	return fmt.Sprintf("%05d", count+1), nil
}
