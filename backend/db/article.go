package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Article struct {
	ID          int
	RssID       int
	Title       string
	Link        string
	GUID        string
	Description string
	PublishDate string
	Format      string
	Identifier  string
	Read        bool
	Author      string    // RSS feed title as author
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func CreateArticleTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS article (
	id SERIAL PRIMARY KEY,
	rssID INT NOT NULL,
	title TEXT,
	link TEXT,
	GUID TEXT,
	description TEXT,
	publishDate TEXT,
	format TEXT,
	identifier TEXT,
	read BOOLEAN,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT fk_articles_rss 
		FOREIGN KEY (rssID) 
		REFERENCES rss(id) 
		ON DELETE CASCADE 
		ON UPDATE CASCADE,

	CONSTRAINT unique_article_rss_link UNIQUE (rssID, link)
	)`

	_, err := DB.Exec(context.Background(), query)
	return err
}

func CreateArticle(rssID int, title, link, guid, description string, publishDate string, format, identifier string, read bool) (*Article, error) {
	query := `
	INSERT INTO article (rssID, title, link, GUID, description, publishDate, format, identifier, read)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (rssID, link) DO NOTHING
	RETURNING id, rssID, title, link, GUID, description, publishDate, format, identifier, read, created_at, updated_at
	`
	article := &Article{}
	err := DB.QueryRow(context.Background(), query, rssID, title, link, guid, description, publishDate, format, identifier, read).Scan(
		&article.ID, &article.RssID, &article.Title, &article.Link, &article.GUID, &article.Description, &article.PublishDate, &article.Format, &article.Identifier, &article.Read, &article.CreatedAt, &article.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		// Article already existed and wasn't inserted
		return nil, nil // or return a specific "already exists" indicator
	}
	return article, err
}

func GetArticleByRSSID(id, limit int) ([]Article, error) {
	query := `
	SELECT a.id, a.rssID, a.title, a.link, a.GUID, a.description, a.publishDate, a.format, a.identifier, a.read, r.title as author, a.created_at, a.updated_at
	FROM article a
	JOIN rss r ON a.rssID = r.id
	WHERE a.rssid = $1
	AND a.publishDate != '' AND a.publishDate IS NOT NULL
	ORDER BY a.publishDate::TIMESTAMP DESC
	LIMIT $2
	`

	rows, err := DB.Query(context.Background(), query, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.RssID, &article.Title, &article.Link,
			&article.GUID, &article.Description, &article.PublishDate, &article.Format, &article.Identifier,
			&article.Read, &article.Author, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

func GetSingleArticle(id int) ([]Article, error) {
	query := `
	SELECT id, rssID, title, link, GUID, description, publishDate, format, identifier, read, created_at, updated_at
	FROM article
	WHERE id = $1
	`

	rows, err := DB.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.RssID, &article.Title, &article.Link,
			&article.GUID, &article.Description, &article.PublishDate, &article.Format, &article.Identifier,
			&article.Read, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

func UpdateArticleReadStatus(id int, read bool) error {
	query := `
	UPDATE article
	SET read = $1
	WHERE id = $2
	`

	result, err := DB.Exec(context.Background(), query, read, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("RSS with ID %d not found", id)
	}

	return nil
}

func GetAllArticles() ([]Article, error) {
	return GetAllArticlesPaginated(0, 100) // Default to first 100 articles
}

func GetAllArticlesPaginated(offset, limit int) ([]Article, error) {
	query := `
	SELECT a.id, a.rssID, a.title, a.link, a.GUID, a.description, a.publishDate, a.format, a.identifier, a.read, r.title as author, a.created_at, a.updated_at
	FROM article a
	JOIN rss r ON a.rssID = r.id
	WHERE a.publishDate != '' AND a.publishDate IS NOT NULL
	ORDER BY a.publishDate::TIMESTAMP DESC
	LIMIT $1 OFFSET $2;
	`

	rows, err := DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.RssID, &article.Title, &article.Link,
			&article.GUID, &article.Description, &article.PublishDate, &article.Format, &article.Identifier,
			&article.Read, &article.Author, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

func GetTotalArticleCount() (int, error) {
	query := `SELECT COUNT(*) FROM article WHERE publishDate != '' AND publishDate IS NOT NULL`
	
	var count int
	err := DB.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}

func SearchArticles(query string, limit int) ([]Article, error) {
	searchQuery := `
	SELECT id, rssID, title, link, GUID, description, publishDate, format, identifier, read, created_at, updated_at
	FROM article
	WHERE to_tsvector('english', title || ' ' || description) @@ plainto_tsquery('english', $1)
	AND publishDate != '' AND publishDate IS NOT NULL
	ORDER BY publishDate::TIMESTAMP DESC
	LIMIT $2
	`

	rows, err := DB.Query(context.Background(), searchQuery, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.RssID, &article.Title, &article.Link,
			&article.GUID, &article.Description, &article.PublishDate, &article.Format, &article.Identifier,
			&article.Read, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

func DeleteArticle(id int) error {
	query := `
	DELETE FROM article
	WHERE id = $1
	`

	result, err := DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("article with ID %d not found", id)
	}

	return nil
}
