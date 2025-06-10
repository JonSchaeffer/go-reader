package db

import (
	"context"
	"time"
)

type RSS struct {
	ID          int
	URL         string
	Title       string
	Description string
	FeedSize    int
	Sync        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func CreateRSSTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS rss (
	id SERIAL PRIMARY KEY,
	url TEXT,
	title TEXT,
	description TEXT,
	feedSize INT,
	sync INT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(context.Background(), query)
	return err
}

func CreateRSS(url, title, description string, feedSize, sync int) (*RSS, error) {
	query := `
	INSERT INTO rss (url, title, description, feedSize, sync) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, url, title, description, feedSize, sync, created_at, updated_at`

	rss := &RSS{}
	err := DB.QueryRow(context.Background(), query, url, title, description, feedSize, sync).Scan(
		&rss.ID, &rss.URL, &rss.Title, &rss.Description, &rss.FeedSize, &rss.Sync,
		&rss.CreatedAt, &rss.UpdatedAt,
	)

	return rss, err
}
