package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init(datasource string) error {
	var err error
	DB, err = pgxpool.New(context.Background(), datasource)
	if err != nil {
		return err
	}
	return DB.Ping(context.Background())
}

func Close() {
	DB.Close()
}
