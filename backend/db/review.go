package db

import (
	"context"
	"time"
)

type ReviewHighlight struct {
	HighlightWithContext
	ReviewCount    int        `json:"ReviewCount"`
	LastReviewedAt *time.Time `json:"LastReviewedAt"`
}

// MigrateReviewColumns adds review tracking columns to the highlight table (idempotent).
func MigrateReviewColumns() error {
	_, err := DB.Exec(context.Background(), `
		ALTER TABLE highlight ADD COLUMN IF NOT EXISTS last_reviewed_at TIMESTAMP;
		ALTER TABLE highlight ADD COLUMN IF NOT EXISTS review_count INT NOT NULL DEFAULT 0;
	`)
	return err
}

// GetReviewQueue returns up to `limit` highlights not yet reviewed today,
// ordered: never-reviewed first, then oldest last_reviewed_at.
func GetReviewQueue(limit int) ([]ReviewHighlight, error) {
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
		) AS item_url,
		h.review_count,
		h.last_reviewed_at
	FROM highlight h
	LEFT JOIN saved_item si ON h.item_type = 'library' AND h.item_id = si.id
	LEFT JOIN article   a  ON h.item_type = 'article'  AND h.item_id = a.id
	WHERE h.last_reviewed_at IS NULL
	   OR h.last_reviewed_at < CURRENT_DATE
	ORDER BY
		h.last_reviewed_at ASC NULLS FIRST,
		h.created_at ASC
	LIMIT $1`

	rows, err := DB.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ReviewHighlight
	for rows.Next() {
		var r ReviewHighlight
		err := rows.Scan(
			&r.ID, &r.ItemType, &r.ItemID, &r.SelectedText, &r.Color, &r.Note, &r.TextOffset, &r.CreatedAt,
			&r.ItemTitle, &r.ItemURL,
			&r.ReviewCount, &r.LastReviewedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	if results == nil {
		results = []ReviewHighlight{}
	}
	return results, rows.Err()
}

// CompleteReview stamps a highlight as reviewed now and increments the count.
func CompleteReview(id int) error {
	_, err := DB.Exec(context.Background(),
		`UPDATE highlight SET last_reviewed_at = NOW(), review_count = review_count + 1 WHERE id = $1`,
		id,
	)
	return err
}

// GetReviewQueueSize returns how many highlights are waiting to be reviewed today.
func GetReviewQueueSize() (int, error) {
	var count int
	err := DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM highlight
		WHERE last_reviewed_at IS NULL OR last_reviewed_at < CURRENT_DATE
	`).Scan(&count)
	return count, err
}
