package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type RSS struct {
	ID          int       `json:"ID"`
	URL         string    `json:"Url"`
	FiveURL     string    `json:"FivefiltersUrl"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	FeedSize    int       `json:"FeedSize"`
	Sync        int       `json:"Sync"`
	CategoryID  *int      `json:"CategoryID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

type Category struct {
	ID        int
	Name      string
	Color     string // Hex color for visual distinction
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateCategoryTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS category (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	color TEXT DEFAULT '#3b82f6',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT unique_category_name UNIQUE (name)
	)`
	_, err := DB.Exec(context.Background(), query)
	return err
}

func CreateRSSTable() error {
	// First ensure category table exists
	if err := CreateCategoryTable(); err != nil {
		return err
	}

	// Create table without foreign key constraint first
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
	if err != nil {
		return err
	}

	return nil
}

func CreateRSS(url, fiveURL, title, description string, feedSize, sync int) (*RSS, error) {
	query := `
	INSERT INTO rss (url, fiveURL, title, description, feedSize, sync, categoryID) 
	VALUES ($1, $2, $3, $4, $5, $6, NULL)
	ON CONFLICT (url) DO NOTHING
	RETURNING id, url, fiveURL, title, description, feedSize, sync, categoryID, created_at, updated_at`

	rss := &RSS{}
	err := DB.QueryRow(context.Background(), query, url, fiveURL, title, description, feedSize, sync).Scan(
		&rss.ID, &rss.URL, &rss.FiveURL, &rss.Title, &rss.Description, &rss.FeedSize, &rss.Sync, &rss.CategoryID,
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
	SELECT id, url, fiveurl, title, description, feedSize, sync, categoryID, created_at, updated_at 
	FROM rss 
	ORDER BY categoryID NULLS FIRST, id`

	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rssFeeds []RSS
	for rows.Next() {
		var rss RSS
		err := rows.Scan(&rss.ID, &rss.URL, &rss.FiveURL, &rss.Title, &rss.Description,
			&rss.FeedSize, &rss.Sync, &rss.CategoryID, &rss.CreatedAt, &rss.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rssFeeds = append(rssFeeds, rss)
	}
	return rssFeeds, rows.Err()
}

func GetRSSByID(id int) (*RSS, error) {
	query := `
	SELECT id, url, fiveurl, title, description, feedSize, sync, categoryID, created_at, updated_at
	FROM rss
	WHERE id = $1
	`

	rss := &RSS{}
	err := DB.QueryRow(context.Background(), query, id).Scan(&rss.ID, &rss.URL, &rss.FiveURL,
		&rss.Title, &rss.Description, &rss.FeedSize, &rss.Sync, &rss.CategoryID, &rss.CreatedAt, &rss.UpdatedAt)

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

// Special function for updating categoryID that handles NULL values
func UpdateRSSCategoryID(id int, categoryID *int) error {
	query := `UPDATE rss SET categoryID = $1 WHERE id = $2`

	result, err := DB.Exec(context.Background(), query, categoryID, id)
	if err != nil {
		log.Printf("Error updating categoryID for RSS %d: %v", id, err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("RSS with ID %d not found", id)
	}

	log.Printf("Successfully updated RSS %d categoryID to %v", id, categoryID)
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

// Category management functions

func CreateCategory(name, color string) (*Category, error) {
	query := `
	INSERT INTO category (name, color) 
	VALUES ($1, $2)
	ON CONFLICT (name) DO NOTHING
	RETURNING id, name, color, created_at, updated_at`

	category := &Category{}
	err := DB.QueryRow(context.Background(), query, name, color).Scan(
		&category.ID, &category.Name, &category.Color, &category.CreatedAt, &category.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("category with name '%s' already exists", name)
	}

	return category, err
}

func GetAllCategories() ([]Category, error) {
	query := `
	SELECT id, name, color, created_at, updated_at 
	FROM category 
	ORDER BY name`

	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Color, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func GetCategoryByID(id int) (*Category, error) {
	query := `
	SELECT id, name, color, created_at, updated_at
	FROM category
	WHERE id = $1
	`

	category := &Category{}
	err := DB.QueryRow(context.Background(), query, id).Scan(&category.ID, &category.Name, &category.Color, &category.CreatedAt, &category.UpdatedAt)

	return category, err
}

func UpdateCategory(id int, name, color string) error {
	query := `UPDATE category SET name = $1, color = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`

	result, err := DB.Exec(context.Background(), query, name, color, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", id)
	}
	return nil
}

func DeleteCategoryByID(id int) error {
	// First, set all RSS feeds in this category to NULL
	_, err := DB.Exec(context.Background(), "UPDATE rss SET categoryID = NULL WHERE categoryID = $1", id)
	if err != nil {
		return err
	}

	// Then delete the category
	query := `DELETE FROM category WHERE id = $1`
	result, err := DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", id)
	}

	return nil
}

func GetRSSByCategory(categoryID *int) ([]RSS, error) {
	var query string
	var args []interface{}

	if categoryID == nil {
		// Get uncategorized feeds
		query = `
		SELECT id, url, fiveurl, title, description, feedSize, sync, categoryID, created_at, updated_at 
		FROM rss 
		WHERE categoryID IS NULL
		ORDER BY id`
		args = []interface{}{}
	} else {
		// Get feeds in specific category
		query = `
		SELECT id, url, fiveurl, title, description, feedSize, sync, categoryID, created_at, updated_at 
		FROM rss 
		WHERE categoryID = $1
		ORDER BY id`
		args = []interface{}{*categoryID}
	}

	rows, err := DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rssFeeds []RSS
	for rows.Next() {
		var rss RSS
		err := rows.Scan(&rss.ID, &rss.URL, &rss.FiveURL, &rss.Title, &rss.Description,
			&rss.FeedSize, &rss.Sync, &rss.CategoryID, &rss.CreatedAt, &rss.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rssFeeds = append(rssFeeds, rss)
	}
	return rssFeeds, rows.Err()
}
