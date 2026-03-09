package db

import (
	"context"
	"fmt"
	"time"
)

type SavedItem struct {
	ID            int       `json:"ID"`
	SourceType    string    `json:"SourceType"` // 'rss_article' | 'url'
	SourceID      *int      `json:"SourceID"`
	URL           string    `json:"URL"`
	Title         string    `json:"Title"`
	Author        string    `json:"Author"`
	Content       string    `json:"Content"`
	PublishedDate string    `json:"PublishedDate"`
	SavedAt       time.Time `json:"SavedAt"`
	Read          bool      `json:"Read"`
	Archived      bool      `json:"Archived"`
}

func CreateSavedItemTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS saved_item (
		id SERIAL PRIMARY KEY,
		source_type TEXT NOT NULL DEFAULT 'url',
		source_id INT REFERENCES article(id) ON DELETE SET NULL,
		url TEXT,
		title TEXT,
		author TEXT,
		content TEXT,
		published_date TEXT,
		saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		read BOOLEAN DEFAULT false,
		archived BOOLEAN DEFAULT false
	)`
	_, err := DB.Exec(context.Background(), query)
	return err
}

func CreateSavedItem(sourceType string, sourceID *int, url, title, author, content, publishedDate string) (*SavedItem, error) {
	query := `
	INSERT INTO saved_item (source_type, source_id, url, title, author, content, published_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, source_type, source_id, url, title, author, content, published_date, saved_at, read, archived`

	item := &SavedItem{}
	err := DB.QueryRow(context.Background(), query, sourceType, sourceID, url, title, author, content, publishedDate).Scan(
		&item.ID, &item.SourceType, &item.SourceID, &item.URL, &item.Title, &item.Author,
		&item.Content, &item.PublishedDate, &item.SavedAt, &item.Read, &item.Archived,
	)
	return item, err
}

func GetAllSavedItems(offset, limit int, archived bool) ([]SavedItem, error) {
	query := `
	SELECT id, source_type, source_id, url, title, author, content, published_date, saved_at, read, archived
	FROM saved_item
	WHERE archived = $1
	ORDER BY saved_at DESC
	LIMIT $2 OFFSET $3`

	rows, err := DB.Query(context.Background(), query, archived, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []SavedItem
	for rows.Next() {
		var item SavedItem
		err := rows.Scan(&item.ID, &item.SourceType, &item.SourceID, &item.URL, &item.Title,
			&item.Author, &item.Content, &item.PublishedDate, &item.SavedAt, &item.Read, &item.Archived)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func GetTotalSavedItemCount(archived bool) (int, error) {
	var count int
	err := DB.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM saved_item WHERE archived = $1", archived).Scan(&count)
	return count, err
}

func GetSavedItemByID(id int) (*SavedItem, error) {
	query := `
	SELECT id, source_type, source_id, url, title, author, content, published_date, saved_at, read, archived
	FROM saved_item WHERE id = $1`

	item := &SavedItem{}
	err := DB.QueryRow(context.Background(), query, id).Scan(
		&item.ID, &item.SourceType, &item.SourceID, &item.URL, &item.Title, &item.Author,
		&item.Content, &item.PublishedDate, &item.SavedAt, &item.Read, &item.Archived,
	)
	return item, err
}

func UpdateSavedItemStatus(id int, read, archived bool) error {
	_, err := DB.Exec(context.Background(),
		"UPDATE saved_item SET read = $1, archived = $2 WHERE id = $3",
		read, archived, id)
	return err
}

func DeleteSavedItem(id int) error {
	result, err := DB.Exec(context.Background(), "DELETE FROM saved_item WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("saved item with ID %d not found", id)
	}
	return nil
}

func SavedItemExistsBySourceID(sourceID int) (bool, error) {
	var exists bool
	err := DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM saved_item WHERE source_id = $1)", sourceID).Scan(&exists)
	return exists, err
}

func SavedItemExistsByURL(url string) (bool, error) {
	var exists bool
	err := DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM saved_item WHERE url = $1)", url).Scan(&exists)
	return exists, err
}

func SearchSavedItems(query string, limit int) ([]SavedItem, error) {
	searchQuery := `
	SELECT id, source_type, source_id, url, title, author, content, published_date, saved_at, read, archived
	FROM saved_item
	WHERE to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content, '')) @@ plainto_tsquery('english', $1)
	ORDER BY saved_at DESC
	LIMIT $2`

	rows, err := DB.Query(context.Background(), searchQuery, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []SavedItem
	for rows.Next() {
		var item SavedItem
		err := rows.Scan(&item.ID, &item.SourceType, &item.SourceID, &item.URL, &item.Title,
			&item.Author, &item.Content, &item.PublishedDate, &item.SavedAt, &item.Read, &item.Archived)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
