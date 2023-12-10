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

	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

// Время ожидания пинга для проверки подключения к базе данных
const pingTimeout = time.Second * 3

// Время ожидания создания таблицы
const createTablesTimeout = time.Second * 5

// DBStorage - структура для взаимодействия с базой данных
type DBStorage struct {
	log *zap.Logger
	db  *sql.DB
}

// ErrConflict - ошибка конфликта данных
var ErrConflict = errors.New("data conflict in DBStorage")

// New - создает новый экземпляр DBStorage
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

// Stop - закрывает соединение с базой данных
func (s *DBStorage) Stop() error {
	s.db.Close()

	return nil
}

// CheckConnection - проверяет подключение к базе данных
func (s *DBStorage) CheckConnection() bool {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	err := s.db.PingContext(ctx)

	if err != nil {
		log.Printf("error pinging the database: %v", err)
	}
	return err == nil
}

// CreateTable - создает таблицу в базе данных
func (s *DBStorage) CreateTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), createTablesTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, createTable)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%v", err)
	}

	return nil
}

// Put - добавляет запись в базу данных
func (s *DBStorage) Put(ctx context.Context, url string, userID int) (string, error) {
	log := lg.New()
	alias := generatestring.NewRandomString(8)

	// Проверяем, является ли строка JSON
	var jsonData map[string]string
	if err := json.Unmarshal([]byte(url), &jsonData); err == nil {
		// Если декодирование прошло успешно, используем значение "url"
		url = jsonData["url"]
	}

	shortURL := ShortURL{
		ID:        userID,
		URL:       url,
		Alias:     alias,
		CreatedAt: time.Now(),
	}

	_, err := s.db.ExecContext(ctx, insert, shortURL.ID, shortURL.URL, shortURL.Alias, shortURL.CreatedAt)
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

// Get - получает URL по алиасу
func (s *DBStorage) Get(ctx context.Context, alias string) (string, error) {
	var url ShortURL
	err := s.db.QueryRowContext(ctx, get, alias).Scan(&url.ID, &url.URL, &url.Alias, &url.CreatedAt)
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

type URLs struct {
	LongLink  string `json:"original_url"`
	ShortLink string `json:"short_url"`
}

func (s *DBStorage) GetURLsByID(ctx context.Context, id int, URLaddr string) ([]byte, error) {

	var userIDs []URLs

	allIDs, err := s.db.QueryContext(ctx, insert, id)
	if err != nil {
		s.log.Error("Error getting batch data: ", zap.Error(err))
		return nil, err
	}

	for allIDs.Next() {
		var links URLs
		err := allIDs.Scan(&links.LongLink, &links.ShortLink)
		if err != nil {
			s.log.Error("Error scanning data: ", zap.Error(err))
			return nil, err
		}
		userIDs = append(userIDs, URLs{
			LongLink:  links.LongLink,
			ShortLink: URLaddr + "/" + links.ShortLink,
		})
	}
	jsonUserIDs, err := json.Marshal(userIDs)
	if err != nil {
		s.log.Error("Can't marshal IDs: ", zap.Error(err))
		return nil, err
	}

	return jsonUserIDs, nil
}

// func (s *DBStorage) GetAll(ctx context.Context, id int, host string) ([]byte, error) {
// 	var userID []ShortURL
// 	allIDs, err := s.db.QueryContext(ctx, getAll, id)
// 	if err != nil {
// 		s.log.Error("Error getting batch data: ", zap.Error(err))
// 		return nil, err
// 	}
// 	defer func() {
// 		_ = allIDs.Close()
// 		_ = allIDs.Err()
// 	}()
//
// 	for allIDs.Next() {
// 		var uid ShortURL
// 		err := allIDs.Scan(&uid.URL, &uid.Alias)
// 		if err != nil {
// 			s.log.Error("Error scanning data: ", zap.Error(err))
// 			return nil, err
// 		}
// 		// Формируем полный URL, включая хост
// 		uid.Alias = host + "/" + uid.Alias
// 		userID = append(userID, uid)
// 		s.log.Info("uid.Alias", zap.String("OUTPUT", uid.Alias))
// 	}
// 	jsonUserIDs, err := json.Marshal(userID)
// 	if err != nil {
// 		s.log.Error("Can't marshal IDs: ", zap.Error(err))
// 		return nil, err
// 	}
// 	return jsonUserIDs, nil
// }
