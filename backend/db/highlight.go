package db

import (
	"context"
	"fmt"
	"time"
)

type Highlight struct {
	ID           int       `json:"ID"`
	ItemType     string    `json:"ItemType"`     // 'library' | 'article'
	ItemID       int       `json:"ItemID"`
	SelectedText string    `json:"SelectedText"`
	Color        string    `json:"Color"` // 'yellow' | 'green' | 'blue' | 'pink'
	Note         string    `json:"Note"`
	TextOffset   int       `json:"TextOffset"` // char offset in article plain text (-1 = unknown)
	CreatedAt    time.Time `json:"CreatedAt"`
}

func CreateHighlightTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS highlight (
		id SERIAL PRIMARY KEY,
		item_type TEXT NOT NULL,
		item_id INT NOT NULL,
		selected_text TEXT NOT NULL,
		color TEXT NOT NULL DEFAULT 'yellow',
		note TEXT NOT NULL DEFAULT '',
		text_offset INT NOT NULL DEFAULT -1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err := DB.Exec(context.Background(), query); err != nil {
		return err
	}
	// Add column to existing tables (idempotent)
	_, err := DB.Exec(context.Background(),
		`ALTER TABLE highlight ADD COLUMN IF NOT EXISTS text_offset INT NOT NULL DEFAULT -1`)
	return err
}

func CreateHighlight(itemType string, itemID int, selectedText, color, note string, textOffset int) (*Highlight, error) {
	query := `
	INSERT INTO highlight (item_type, item_id, selected_text, color, note, text_offset)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, item_type, item_id, selected_text, color, note, text_offset, created_at`

	h := &Highlight{}
	err := DB.QueryRow(context.Background(), query, itemType, itemID, selectedText, color, note, textOffset).Scan(
		&h.ID, &h.ItemType, &h.ItemID, &h.SelectedText, &h.Color, &h.Note, &h.TextOffset, &h.CreatedAt,
	)
	return h, err
}

func GetHighlightsByItem(itemType string, itemID int) ([]Highlight, error) {
	query := `
	SELECT id, item_type, item_id, selected_text, color, note, text_offset, created_at
	FROM highlight
	WHERE item_type = $1 AND item_id = $2
	ORDER BY created_at ASC`

	rows, err := DB.Query(context.Background(), query, itemType, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var highlights []Highlight
	for rows.Next() {
		var h Highlight
		err := rows.Scan(&h.ID, &h.ItemType, &h.ItemID, &h.SelectedText, &h.Color, &h.Note, &h.TextOffset, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		highlights = append(highlights, h)
	}
	return highlights, rows.Err()
}

func UpdateHighlight(id int, color, note string) error {
	result, err := DB.Exec(context.Background(),
		"UPDATE highlight SET color = $1, note = $2 WHERE id = $3",
		color, note, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("highlight %d not found", id)
	}
	return nil
}

type HighlightWithContext struct {
	Highlight
	ItemTitle string `json:"ItemTitle"`
	ItemURL   string `json:"ItemURL"`
}

func GetAllHighlights() ([]HighlightWithContext, error) {
	query := `
	SELECT h.id, h.item_type, h.item_id, h.selected_text, h.color, h.note, h.text_offset, h.created_at,
		COALESCE(
			CASE WHEN h.item_type = 'library' THEN si.title
			     WHEN h.item_type = 'article' THEN a.title
			END, ''
		) AS item_title,
		COALESCE(
			CASE WHEN h.item_type = 'library' THEN si.url
			     WHEN h.item_type = 'article' THEN a.link
			END, ''
		) AS item_url
	FROM highlight h
	LEFT JOIN saved_item si ON h.item_type = 'library' AND h.item_id = si.id
	LEFT JOIN article a ON h.item_type = 'article' AND h.item_id = a.id
	ORDER BY h.created_at DESC`

	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []HighlightWithContext
	for rows.Next() {
		var h HighlightWithContext
		err := rows.Scan(
			&h.ID, &h.ItemType, &h.ItemID, &h.SelectedText, &h.Color, &h.Note, &h.TextOffset, &h.CreatedAt,
			&h.ItemTitle, &h.ItemURL,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, h)
	}
	return results, rows.Err()
}

func DeleteHighlight(id int) error {
	result, err := DB.Exec(context.Background(), "DELETE FROM highlight WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("highlight %d not found", id)
	}
	return nil
}
