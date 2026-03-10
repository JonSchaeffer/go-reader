package db

import (
	"context"
)

type SearchResult struct {
	ID          int    `json:"ID"`
	Type        string `json:"Type"`   // 'article' | 'library'
	Title       string `json:"Title"`
	Excerpt     string `json:"Excerpt"`
	Author      string `json:"Author"`
	PublishDate string `json:"PublishDate"`
	URL         string `json:"URL"`
	Read        bool   `json:"Read"`
}

// SearchAll runs full-text search across both articles and saved_items,
// returning a unified result list ordered by relevance.
func SearchAll(query string, limit int) ([]SearchResult, error) {
	sql := `
	SELECT id, 'article' AS type, title,
		left(regexp_replace(description, '<[^>]*>', '', 'g'), 200) AS excerpt,
		'' AS author, publishdate AS publish_date, link AS url, read
	FROM article
	WHERE to_tsvector('english', coalesce(title,'') || ' ' || coalesce(description,''))
	      @@ plainto_tsquery('english', $1)

	UNION ALL

	SELECT id, 'library' AS type, title,
		left(regexp_replace(content, '<[^>]*>', '', 'g'), 200) AS excerpt,
		author, published_date AS publish_date, url, read
	FROM saved_item
	WHERE archived = false
	  AND to_tsvector('english', coalesce(title,'') || ' ' || coalesce(content,''))
	      @@ plainto_tsquery('english', $1)

	ORDER BY title
	LIMIT $2`

	rows, err := DB.Query(context.Background(), sql, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		if err := rows.Scan(&r.ID, &r.Type, &r.Title, &r.Excerpt, &r.Author, &r.PublishDate, &r.URL, &r.Read); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
