package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/entity"
	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
)

const (
	// createTable - creating a table
	createTable = `CREATE TABLE IF NOT EXISTS short_urls (
    uuid INT,
    url VARCHAR NOT NULL,
    alias VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    del BOOLEAN,
    PRIMARY KEY (uuid, alias),
    UNIQUE (uuid, url)
);`
	createTablesTimeout = time.Second * 5

	// insert SQL query to insert a new short URL record into the short_urls table
	insert = `INSERT INTO short_urls (uuid, url, alias, created_at, del) VALUES ($1, $2, $3, $4, false);`

	// get SQL query to retrieve a short URL record by alias from the short_urls table
	get = `SELECT uuid, url, alias, created_at, del FROM short_urls WHERE alias = $1;`

	// getConflict SQL query to check for conflicts by retrieving alias for a given URL from the short_urls table
	getConflict = `SELECT alias FROM short_urls WHERE url = $1;`
)

// ErrConflict - data conflict in data base storage
var ErrConflict = errors.New("data conflict in DBStorage")

// NewDB - creates a new DBStorage instance
func NewDB(cfg string, log *zap.Logger) (*UseCase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), createTablesTimeout)
	defer cancel()
	// Creating a Database Connection Using Context.
	DB, err := sql.Open("pgx", cfg)
	if err != nil {
		log.Error("error when opening a connection to the database", zap.Error(err))
		return nil, fmt.Errorf("DB connection err=%w", err)
	}
	// Testing database connection using context.
	if err := DB.PingContext(ctx); err != nil {
		log.Error("error when checking database connection", zap.Error(err))
		return nil, fmt.Errorf("DB ping error: %w", err)
	}
	storage := &UseCase{
		DB: DB,
	}
	// create a table using context.
	if err := storage.CreateTable(ctx); err != nil {
		return nil, fmt.Errorf("create table error: %w", err)
	}
	return storage, nil
}

// pingTimeout to check database connection.
const pingTimeout = time.Second * 3

// Stop - closes the connection to the database.
func (uc *UseCase) Stop() error {
	uc.DB.Close()
	return nil
}

// Healthcheck - checks the connection to the database.
func (uc *UseCase) Healthcheck() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := uc.DB.PingContext(ctx); err != nil {
		// Handling a connection error, for example logging or returning false.
		uc.log.Error("Ошибка подключения к базе данных", zap.Error(err))
		return false, err
	}
	return true, nil
}

// createTable - creates a table in the database.
func (uc *UseCase) CreateTable(ctx context.Context) error {
	_, err := uc.DB.ExecContext(ctx, createTable)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%v", err)
	}

	return nil
}

// Put - adds an entry to the database.
func (uc *UseCase) Put(ctx context.Context, url string, uuid int) (string, error) {
	alias := generatestring.NewRandomString(8)

	var jsonData map[string]string
	if err := json.Unmarshal([]byte(url), &jsonData); err == nil {
		// If decoding was successful, use the value "url".
		url = jsonData["url"]
	}

	shortURL := entity.DBStorage{
		UUID:      uuid,
		URL:       url,
		Alias:     alias,
		CreatedAt: time.Now(),
	}

	_, err := uc.DB.ExecContext(ctx, insert, shortURL.UUID, shortURL.URL, shortURL.Alias, shortURL.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			// In case of a conflict, we perform an additional request to obtain the alias.
			var existingAlias string
			err := uc.DB.QueryRowContext(ctx, getConflict, url).Scan(&existingAlias)
			if err != nil {
				if err == sql.ErrNoRows {
					// Handling the situation when a URL match is not found.
					return alias, ErrConflict
				}
				return alias, fmt.Errorf("failed to query existing alias: %w", err)
			}
			return existingAlias, ErrConflict
		}
		return alias, fmt.Errorf("failed to insert short URL into database: %w", err)
	}
	return alias, nil
}

// Get - gets a URL by alias.
func (uc *UseCase) Get(ctx context.Context, alias string) (string, bool, error) {
	var url entity.DBStorage
	err := uc.DB.QueryRowContext(ctx, get, alias).Scan(&url.UUID, &url.URL, &url.Alias, &url.CreatedAt, &url.DeletedFlag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// This is an expected error when there are no rows matching the query.
			return "", false, fmt.Errorf("no URL found for alias %s", alias)
		}
		// Handling other database errors.
		return "", false, err
	}
	return url.URL, url.DeletedFlag, nil
}

// GetAll - Gets all URLs for a specific user.
func (uc *UseCase) GetAll(ctx context.Context, id int, host string) ([]byte, error) {
	var urls []entity.DBStorage
	DB := bun.NewDB(uc.DB, pgdialect.New())

	rows, err := DB.NewSelect().
		TableExpr("short_urls").
		Column("url", "alias").
		Where("uuid = ?", id).
		Rows(ctx)
	if err != nil {
		uc.log.Error("Error getting data: ", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url entity.DBStorage
		if err := rows.Scan(&url.URL, &url.Alias); err != nil {
			uc.log.Error("Error scanning data: ", zap.Error(err))
			return nil, err
		}
		url.Alias = fmt.Sprintf("%s/%s", host, url.Alias)
		urls = append(urls, url)
	}

	if rows.Err() != nil {
		uc.log.Error("Error iterating over rows: ", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	allURL, err := json.Marshal(urls)
	if err != nil {
		uc.log.Error("Can't marshal URLs: ", zap.Error(err))
		return nil, err
	}

	return allURL, nil
}

// Del removes URLs for users with a specific ID.
func (uc *UseCase) Del(ctx context.Context, id int, aliases []string) error {
	inputCh := delGenerator(ctx, aliases)
	if err := uc.updateStatusDel(ctx, id, inputCh); err != nil {
		return fmt.Errorf("failed to delete URLs: %w", err)
	}
	return nil
}

// delGenerator - channel generator for collecting aliases.
func delGenerator(ctx context.Context, URLs []string) chan string {
	URLCh := make(chan string)
	go func() {
		defer close(URLCh)
		for _, data := range URLs {
			select {
			case <-ctx.Done():
				return
			case URLCh <- data:
			}
		}
	}()
	return URLCh
}

// updateStatusDel - performs a batch update of the deletion status for the specified aliases and user.
func (uc *UseCase) updateStatusDel(ctx context.Context, id int, inputCh chan string) error {
	var deleteURLs []string

	// a channel for signaling the end of each goroutine.
	aliasCollectionDoneCh := make(chan struct{})

	// launching a goroutine to collect aliases
	go func() {
		defer close(aliasCollectionDoneCh)
		for alias := range inputCh {
			deleteURLs = append(deleteURLs, alias)
		}
	}()

	// waiting for the goroutine to finish collecting aliases.
	<-aliasCollectionDoneCh

	DB := bun.NewDB(uc.DB, pgdialect.New())

	_, err := DB.NewUpdate().
		TableExpr("short_urls").
		Set("del = ?", true).
		Where("alias IN (?)", bun.In(deleteURLs)).
		Where("uuid = ?", id).
		Exec(ctx)

	if err != nil {
		uc.log.Error("Can't exec update request: ", zap.Error(err))
		return fmt.Errorf("failed to update URLs: %w", err)
	}

	return nil
}
