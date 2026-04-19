package postgres_pool

import (
	"context"
	"fmt"

	"github.com/chixxx1/subscription-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool struct {
	*pgxpool.Pool
}

func NewConnectionPool(ctx context.Context, ctf config.DBConfig) (*ConnectionPool, error) {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", ctf.User, ctf.Password, ctf.Host, ctf.Port, ctf.Database)

	pgxconfig, err := pgxpool.ParseConfig(connectionStr)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &ConnectionPool{
		Pool: pool,
	}, nil
}
