//go:build integration
// +build integration

package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"timescale/internal/config"
)

type PGClientLocal struct {
	db *pgxpool.Pool
}

func New(conf config.Config) (*PGClientLocal, error) {
	conn, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.Postgres.Username, conf.Postgres.Password, conf.Postgres.Hostname, conf.Postgres.Port, conf.Postgres.Database))
	if err != nil {
		return nil, err
	}

	return &PGClientLocal{
		db: conn,
	}, nil
}

func (p PGClientLocal) Create(ctx context.Context, ts time.Time, host string, usage float64) error {

	createQuery := `INSERT INTO cpu_usage (ts, host, usage) VALUES ($1, $2, $3);`
	// Execute INSERT command
	_, err := p.db.Exec(ctx, createQuery, ts, host, usage)
	if err != nil {
		return err
	}
	return nil
}

func (p PGClientLocal) DeleteAll(ctx context.Context) error {

	deleteQuery := `DELETE FROM cpu_usage;`
	// Execute INSERT command
	_, err := p.db.Exec(ctx, deleteQuery)
	if err != nil {
		return err
	}
	return nil
}

func (p PGClientLocal) Close(ctx context.Context) {
	p.db.Close()
}
