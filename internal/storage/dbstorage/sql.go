package dbstorage

import (
	"time"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS short_urls (
    id INT,
    url VARCHAR(255) UNIQUE NOT NULL, 
    alias VARCHAR(255) UNIQUE NOT NULL,
    created_at timestamp
);`

	insert      = `INSERT INTO short_urls (id, url, alias, created_at) VALUES ($1, $2, $3, $4)`
	get         = `SELECT id, url, alias, created_at FROM short_urls WHERE alias = $1;`
	getConflict = `SELECT alias FROM short_urls WHERE url = $1`
	getAll      = `SELECT id, url, alias, created_at FROM short_urls WHERE id = $1;
`
	// update = `UPDATE short_urls SET url=$1, alias=$2, created_at=$3 WHERE id=$4;`
	// delete = `DELETE FROM short_urls WHERE id=$1;`
)

type ShortURL struct {
	ID        int       `json:"user_id,omitempty"`
	URL       string    `json:"original_url,omitempty"`
	Alias     string    `json:"short_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
