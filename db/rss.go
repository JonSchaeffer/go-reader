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

func UpdateRSS[T string | int](id int, param string, value T) error {
	var query string

	switch param {
	case "url":
		query = `UPDATE rss SET url = $1 WHERE id = $2`
	case "feedsize":
		query = `UPDATE rss SET feedsize = $1 WHERE id = $2`
	case "sync":
		query = `UPDATE rss SET sync = $1 WHERE id = $2`
	case "title":
		query = `UPDATE rss SET title = $1 WHERE id = $2`
	case "description":
		query = `UPDATE rss SET description = $1 WHERE id = $2`
	default:
		return fmt.Errorf("invalid parameter: %s", param)
	}

	result, err := DB.Exec(context.Background(), query, value, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("RSS with ID %d not found", id)
	}
	return nil
}

type RSSStats struct {
	FeedID            int       `json:"feed_id"`
	TotalArticles     int       `json:"total_articles"`
	UnreadArticles    int       `json:"unread_articles"`
	ReadArticles      int       `json:"read_articles"`
	OldestArticle     time.Time `json:"oldest_article"`
	NewestArticle     time.Time `json:"newest_article"`
	LastUpdated       time.Time `json:"last_updated"`
	DaysSinceLastPost int       `json:"days_since_last_post"`
}

func GetRSSStats(id int) (*RSSStats, error) {
	stats := &RSSStats{FeedID: id}
	
	// Get total article count
	err := DB.QueryRow(context.Background(), 
		"SELECT COUNT(*) FROM article WHERE rssid = $1", id).Scan(&stats.TotalArticles)
	if err != nil {
		return nil, err
	}
	
	// Get unread article count
	err = DB.QueryRow(context.Background(), 
		"SELECT COUNT(*) FROM article WHERE rssid = $1 AND read = false", id).Scan(&stats.UnreadArticles)
	if err != nil {
		return nil, err
	}
	
	// Calculate read articles
	stats.ReadArticles = stats.TotalArticles - stats.UnreadArticles
	
	// Get oldest and newest article dates (return early if no articles)
	if stats.TotalArticles > 0 {
		err = DB.QueryRow(context.Background(), 
			"SELECT MIN(created_at), MAX(created_at) FROM article WHERE rssid = $1", 
			id).Scan(&stats.OldestArticle, &stats.NewestArticle)
		if err != nil {
			return nil, err
		}
		
		// Get days since last post (using publishDate if available, otherwise created_at)
		var lastPostDate time.Time
		err = DB.QueryRow(context.Background(), `
			SELECT COALESCE(MAX(
				CASE 
					WHEN publishdate != '' AND publishdate IS NOT NULL 
					THEN publishdate::timestamp 
					ELSE created_at 
				END
			), MIN(created_at))
			FROM article WHERE rssid = $1`, id).Scan(&lastPostDate)
		if err != nil {
			return nil, err
		}
		
		stats.DaysSinceLastPost = int(time.Since(lastPostDate).Hours() / 24)
	}
	
	// Get RSS feed last updated time
	err = DB.QueryRow(context.Background(), 
		"SELECT updated_at FROM rss WHERE id = $1", id).Scan(&stats.LastUpdated)
	if err != nil {
		return nil, err
	}
	
	return stats, nil
}
