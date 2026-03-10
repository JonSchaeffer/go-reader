package db

import (
	"context"
)

type Tag struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

func CreateTagTables() error {
	_, err := DB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS tag (
			id   SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS item_tag (
			item_type TEXT NOT NULL,
			item_id   INT  NOT NULL,
			tag_id    INT  NOT NULL REFERENCES tag(id) ON DELETE CASCADE,
			PRIMARY KEY (item_type, item_id, tag_id)
		);
	`)
	return err
}

func GetAllTags() ([]Tag, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, name FROM tag ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, nil
}

func CreateTag(name string) (Tag, error) {
	var t Tag
	err := DB.QueryRow(context.Background(),
		`INSERT INTO tag (name) VALUES ($1)
		 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name`,
		name,
	).Scan(&t.ID, &t.Name)
	return t, err
}

func DeleteTag(id int) error {
	_, err := DB.Exec(context.Background(), `DELETE FROM tag WHERE id = $1`, id)
	return err
}

func GetTagsForItem(itemType string, itemID int) ([]Tag, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT t.id, t.name FROM tag t
		 JOIN item_tag it ON it.tag_id = t.id
		 WHERE it.item_type = $1 AND it.item_id = $2
		 ORDER BY t.name`,
		itemType, itemID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, nil
}

func AddTagToItem(itemType string, itemID, tagID int) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO item_tag (item_type, item_id, tag_id)
		 VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		itemType, itemID, tagID,
	)
	return err
}

func RemoveTagFromItem(itemType string, itemID, tagID int) error {
	_, err := DB.Exec(context.Background(),
		`DELETE FROM item_tag WHERE item_type = $1 AND item_id = $2 AND tag_id = $3`,
		itemType, itemID, tagID,
	)
	return err
}
