package dbstorage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase"
	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

const (
	pingTimeout         = time.Second * 3
	createTablesTimeout = time.Second * 5
)

type DBStorage struct {
	log *zap.Logger
	db  *sql.DB
}

var ErrConflict = errors.New("data conflict in DBStorage")

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

func (s *DBStorage) Put(ctx context.Context, url string, userID int) (string, error) {
	log := lg.New()
	alias := generatestring.NewRandomString(8)
	var shortURL usecase.CustomRequest
	shortURL.URL = url
	shortURL.Alias = alias
	shortURL.UserID = userID
	shortURL.CreatedAt = time.Now()

	_, err := s.db.ExecContext(ctx, insert, shortURL.UserID, shortURL.URL, shortURL.Alias, shortURL.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			// В случае конфликта выполните дополнительный запрос для получения алиаса
			var existingAlias string
			err := s.db.QueryRowContext(ctx, getConflict, url).Scan(&existingAlias)
			if err != nil {
				if err == sql.ErrNoRows {
					// Обработка ситуации, когда не найдено совпадение по URL
					return alias, ErrConflict
				}
				log.Error("DBStorage.QueryRowContext failed", zap.Error(err))
				return alias, fmt.Errorf("failed to query existing alias: %w", err)
			}
			return existingAlias, ErrConflict
		}

		// Логирование текста ошибки для анализа
		log.Error("DBStorage.Put failed", zap.Error(err))
		return alias, fmt.Errorf("failed to insert short URL into database: %w", err)
	}

	// Логирование успешной вставки
	lg.New().Info("DBStorage.Put", zap.String("alias", alias), zap.String("url", url))
	return alias, nil
}

func (s *DBStorage) Get(ctx context.Context, alias string) (string, error) {
	var url usecase.CustomRequest
	err := s.db.QueryRowContext(ctx, get, alias).Scan(&url.UserID, &url.URL, &url.Alias, &url.CreatedAt)
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

func (s *DBStorage) GetAll(ctx context.Context, userID int, url string) ([]byte, error) {
	var uid usecase.CustomRequest
	allURL, err := s.db.QueryContext(ctx, getAll, userID)
	if err != nil {
		s.log.Error("Error getting batch data: ", zap.Error(err))
		return nil, err
	}
	defer func() {
		_ = allURL.Close()
		_ = allURL.Err()
	}()

	var userURL []usecase.CustomRequest
	for allURL.Next() {
		err := allURL.Scan(&uid.UserID, &uid.URL, &uid.Alias, &uid.CreatedAt)
		if err != nil {
			s.log.Error("Error scanning data: ", zap.Error(err))
			return nil, err
		}
		userURL = append(userURL, uid)
	}

	return json.Marshal(userURL)
}
