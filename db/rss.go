package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type RSS struct {
	ID          int
	URL         string
	FiveURL     string
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
	fiveURL TEXT,
	title TEXT,
	description TEXT,
	feedSize INT,
	sync INT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT unique_rss_url UNIQUE (url)
	)`
	_, err := DB.Exec(context.Background(), query)
	return err
}

func CreateRSS(url, fiveURL, title, description string, feedSize, sync int) (*RSS, error) {
	query := `
	INSERT INTO rss (url, fiveURL, title, description, feedSize, sync) 
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (url) DO NOTHING
	RETURNING id, url, fiveURL, title, description, feedSize, sync, created_at, updated_at`

	rss := &RSS{}
	err := DB.QueryRow(context.Background(), query, url, fiveURL, title, description, feedSize, sync).Scan(
		&rss.ID, &rss.URL, &rss.FiveURL, &rss.Title, &rss.Description, &rss.FeedSize, &rss.Sync,
		&rss.CreatedAt, &rss.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		// Article already existed and wasn't inserted
		return nil, nil // or return a specific "already exists" indicator
	}

	return rss, err
}

func GetAllRSS() ([]RSS, error) {
	query := `
	SELECT id, url, fiveurl, title, description, feedSize, sync, created_at, updated_at 
	FROM rss 
	ORDER BY id`

	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rssFeeds []RSS
	for rows.Next() {
		var rss RSS
		err := rows.Scan(&rss.ID, &rss.URL, &rss.FiveURL, &rss.Title, &rss.Description,
			&rss.FeedSize, &rss.Sync, &rss.CreatedAt, &rss.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rssFeeds = append(rssFeeds, rss)
	}
	return rssFeeds, rows.Err()
}

func GetRSSByID(id int) (*RSS, error) {
	query := `
	SELECT id, url, fiveurl, title, description, feedSize, sync, created_at, updated_at
	FROM rss
	WHERE id = $1
	`

	rss := &RSS{}
	err := DB.QueryRow(context.Background(), query, id).Scan(&rss.ID, &rss.URL, &rss.FiveURL,
		&rss.Title, &rss.Description, &rss.FeedSize, &rss.Sync, &rss.CreatedAt, &rss.UpdatedAt)

	return rss, err
}

func DeleteRSSByID(id int) error {
	query := `
	DELETE FROM rss WHERE id = $1
	`

	result, err := DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("RSS with ID %d not found", id)
	}

	// Reset the sequence after successful deletion
	resetQuery := `SELECT setval('rss_id_seq', COALESCE(MAX(id), 0) + 1, false) FROM rss`
	_, err = DB.Exec(context.Background(), resetQuery)
	if err != nil {
		// Log the error but don't fail the function since the delete succeeded
		log.Printf("Warning: failed to reset rss sequence: %v", err)
	}

	return nil
}
