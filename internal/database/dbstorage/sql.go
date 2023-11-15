package dbstorage

import (
	"time"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS short_urls (
    id serial PRIMARY KEY,
    url text NOT NULL,
    alias text NOT NULL UNIQUE,
    created_at timestamp
);`

	insert = `INSERT INTO short_urls (url, alias, created_at) VALUES ($1, $2, $3)  RETURNING id;`
	get    = `SELECT id, url, alias, created_at FROM short_urls WHERE alias = $1;`
	// getAll = `SELECT id, url, alias, created_at FROM short_urls;`
	// update = `UPDATE short_urls SET url=$1, alias=$2, created_at=$3 WHERE id=$4;`
	// delete = `DELETE FROM short_urls WHERE id=$1;`
)

type ShortURL struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
}
