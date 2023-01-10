package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"timescale/internal/config"
)

type Row struct {
	Ts                 time.Time
	MaxUsage, MinUsage float32
}

// Why use github.com/jackc/pgx/v5 ?
// https://levelup.gitconnected.com/fastest-postgresql-client-library-for-go-579fa97909fb
type PGClient struct {
	db *pgxpool.Pool
}

func New(conf config.Config) (*PGClient, error) {
	conn, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.Postgres.Username, conf.Postgres.Password, conf.Postgres.Hostname, conf.Postgres.Port, conf.Postgres.Database))
	if err != nil {
		return nil, err
	}

	return &PGClient{
		db: conn,
	}, nil
}

func (p PGClient) GetStatsPerHostPerMin(ctx context.Context, host string, start, end time.Time) ([]Row, error) {

	var query = `
				select 
				time_bucket('1 minutes', ts) AS min_interval ,
				max(usage),
				min(usage)
				from cpu_usage
				where host=$1 and 
				ts >= $2 and
				ts <= $3
				GROUP BY min_interval
				ORDER BY min_interval`

	rows, err := p.db.Query(ctx, query, host, start, end)
	if err != nil {
		return nil, err
	}

	var results []Row
	for rows.Next() {
		r := Row{}
		if err := rows.Scan(&r.Ts, &r.MaxUsage, &r.MinUsage); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (p PGClient) Close(ctx context.Context) {
	p.db.Close()
}
