package dbstorage

const (
	createTable = `CREATE TABLE IF NOT EXISTS short_urls (
    uuid INT,
    url VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    del BOOLEAN,
    PRIMARY KEY (uuid, alias),
    UNIQUE (uuid, url)
);`

	insert      = `INSERT INTO short_urls (uuid, url, alias, created_at) VALUES ($1, $2, $3, $4);`
	get         = `SELECT uuid, url, alias, created_at FROM short_urls WHERE alias = $1;`
	getConflict = `SELECT alias FROM short_urls WHERE url = $1;`
	// getAll      = `SELECT url, alias FROM short_urls WHERE uuid = $1;`
	// update = `UPDATE short_urls SET url=$1, alias=$2, created_at=$3 WHERE id=$4;`
	// delete = `DELETE FROM short_urls WHERE id=$1;`
)
