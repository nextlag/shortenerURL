package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
)

const (
	pingTimeout         = time.Second * 3
	createTablesTimeout = time.Second * 5
)

type DBStorage struct {
	db *sql.DB
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

// Put сохраняет значение по ключу
func (s *DBStorage) Put(url string) (string, error) {
	alias := generatestring.NewRandomString(8)
	var id int64
	shortURL := ShortURL{
		URL:       url,
		Alias:     alias,
		CreatedAt: time.Now(),
	}

	err := s.db.QueryRow(insert, shortURL.URL, shortURL.Alias, shortURL.CreatedAt).Scan(&id)
	if err != nil {
		// Если произошел конфликт (например, дублирование URL), возвращаем существующий alias
		if pgerrcode.IsIntegrityConstraintViolation(err.(*pgconn.PgError).Code) {
			existingAlias, err := s.Get(url)
			if err != nil {
				return "", fmt.Errorf("failed to get existing alias: %w", err)
			}
			return existingAlias, nil
		}
		return "", fmt.Errorf("failed to insert short URL into database: %w", err)
	}

	return alias, nil
}

func (s *DBStorage) Get(alias string) (string, error) {
	var url ShortURL
	err := s.db.QueryRow(get, alias).Scan(&url.ID, &url.URL, &url.Alias, &url.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Это ожидаемая ошибка, когда нет строк, соответствующих запросу.
			return "", fmt.Errorf("no URL found for alias %s", alias)
		}
		// Обработка других ошибок базы данных
		return "", err
	}
	return url.URL, nil
}
