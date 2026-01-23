package testdb

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestPostgres struct {
	Container *postgres.PostgresContainer
	DSN       string
}

func StartPostgres(ctx context.Context) (*TestPostgres, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:15",
		postgres.WithDatabase("demobot_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &TestPostgres{
		Container: container,
		DSN:       dsn,
	}, nil
}

func InitPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 5

	return pgxpool.NewWithConfig(ctx, cfg)
}

func ApplyEnv(dsn string) {
	os.Setenv("DB_DSN", dsn)
}
