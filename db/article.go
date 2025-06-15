package db

import (
	"context"
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
