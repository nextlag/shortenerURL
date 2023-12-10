package dbstorage

import (
	"time"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS short_urls (
    user_id INT,
    url VARCHAR(255) UNIQUE NOT NULL, 
    alias VARCHAR(255) UNIQUE NOT NULL,
    created_at timestamp
);`

	insert      = `INSERT INTO short_urls (user_id, url, alias, created_at) VALUES ($1, $2, $3, $4);`
	get         = `SELECT user_id, url, alias, created_at FROM short_urls WHERE alias = $1;`
	getConflict = `SELECT alias FROM short_urls WHERE url = $1;`
	getAll      = `SELECT url, alias FROM short_urls WHERE user_id = $1;`
	// update = `UPDATE short_urls SET url=$1, alias=$2, created_at=$3 WHERE id=$4;`
	// delete = `DELETE FROM short_urls WHERE id=$1;`
)

type ShortURL struct {
	UserID    int       `json:"user_id,omitempty" `
	URL       string    `json:"original_url,omitempty"`
	Alias     string    `json:"short_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
