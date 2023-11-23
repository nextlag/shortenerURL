package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	pingTimeout         = time.Second * 3
	createTablesTimeout = time.Second * 5
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{db: db}
}

func New(dbConfig string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dbConfig)
	if err != nil {
		return nil, fmt.Errorf("db connection err=%w", err)
	}
	storage := &DBStorage{db: db}
	if err := storage.CreateTable(); err != nil {
		return nil, fmt.Errorf("create table error: %w", err)
	}
	return storage, nil
}

func (s *DBStorage) Stop() error {
	s.db.Close()

	return nil
}

func (s *DBStorage) CheckConnection() bool {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	err := s.db.PingContext(ctx)

	if err != nil {
		log.Printf("error pinging the database: %v", err)
	}
	return err == nil
}

func (s *DBStorage) CreateTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), createTablesTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, createTable)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%v", err)
	}

	return nil
}

func (s *DBStorage) Put(alias, url string) error {
	var id int64
	shortURL := &ShortURL{
		URL:       url,
		Alias:     alias,
		CreatedAt: time.Now(),
	}

	err := s.db.QueryRow(insert, shortURL.URL, shortURL.Alias, shortURL.CreatedAt).Scan(&id)
	if err != nil {
		log.Printf("Error inserting short URL into database: %v", err)
		return fmt.Errorf("failed to insert short URL into database: %w", err)
	}
	log.Printf("Inserted short URL with ID: %d", id)
	return nil
}

func (s *DBStorage) Get(alias string) (string, error) {
	var url ShortURL
	err := s.db.QueryRow(get, alias).Scan(&url.URL)
	if err != nil {
		return "", err
	}
	return url.URL, nil
}
