package dbstorage

const (
	createTable = `CREATE TABLE IF NOT EXISTS short_urls (
    userID INT,
    url VARCHAR(255) UNIQUE NOT NULL, 
    alias VARCHAR(255) UNIQUE NOT NULL,
);`

	insert      = `INSERT INTO short_urls (user_id, url, alias) VALUES ($1, $2, $3)`
	get         = `SELECT userID, url, alias, created_at FROM short_urls WHERE alias = $1;`
	getConflict = `SELECT alias FROM short_urls WHERE url = $1`
	getAll      = `SELECT url, alias from short_urls WHERE userID = $1;`
	// update = `UPDATE short_urls SET url=$1, alias=$2, created_at=$3 WHERE id=$4;`
	// delete = `DELETE FROM short_urls WHERE id=$1;`
)
