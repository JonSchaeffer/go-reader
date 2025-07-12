package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init(datasource string) error {
	// Parse the connection string into a config
	config, err := pgxpool.ParseConfig(datasource)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Configure connection pool settings
	config.MaxConns = 10                      // Maximum number of connections in the pool
	config.MinConns = 2                       // Minimum number of connections to maintain
	config.MaxConnLifetime = time.Hour        // Max lifetime of a connection
	config.MaxConnIdleTime = time.Minute * 30 // Max idle time before closing
	config.HealthCheckPeriod = time.Minute    // How often to check connection health

	// Create the connection pool with the configured settings
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.Ping(ctx)
	if err != nil {
		DB.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully with pool config: MaxConns=%d, MinConns=%d",
		config.MaxConns, config.MinConns)

	return nil
}

func Close() {
	DB.Close()
}
