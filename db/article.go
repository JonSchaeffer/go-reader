package db

import "context"

type Article struct{}

func CreateArticleTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS article (
	id SERIAL PRIMARY KEY,
	rssID INT NOT NULL,
	title TEXT,
	link TEXT,
	GUID INT,
	description TEXT,
	publishDate TIMESTAMP,
	format TEXT,
	identifier TEXT,
	read BOOLEAN,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT fk_articles_rss 
		FOREIGN KEY (rssID) 
		REFERENCES rss(id) 
		ON DELETE CASCADE 
		ON UPDATE CASCADE
	)`
	_, err := DB.Exec(context.Background(), query)
	return err
}
