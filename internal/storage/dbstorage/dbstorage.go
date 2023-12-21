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
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
)

type DBStorage struct {
	db          *sql.DB
	log         *zap.Logger
	UUID        int       `json:"user_id,omitempty" `
	URL         string    `json:"original_url,omitempty"`
	Alias       string    `json:"short_url,omitempty"`
	DeletedFlag bool      `json:"deleted_flag"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// Время ожидания пинга для проверки подключения к базе данных
const pingTimeout = time.Second * 3

// ErrConflict - ошибка конфликта данных
var ErrConflict = errors.New("data conflict in DBStorage")

// New - создает новый экземпляр DBStorage
func New(ctx context.Context, cfg string, log *zap.Logger) (*DBStorage, error) {
	// Создание подключения к базе данных с использованием контекста
	db, err := sql.Open("pgx", cfg)
	if err != nil {
		return nil, fmt.Errorf("db connection err=%w", err)
	}
	// Проверка подключения к базе данных с использованием контекста
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}
	storage := &DBStorage{
		db:  db,
		log: log,
	}
	// Создание таблицы с использованием контекста
	if err := storage.CreateTable(ctx); err != nil {
		return nil, fmt.Errorf("create table error: %w", err)
	}
	return storage, nil
}

// Stop - закрывает соединение с базой данных
func (s *DBStorage) Stop() error {
	s.db.Close()

	return nil
}

// Healtcheck - проверяет подключение к базе данных
func (s *DBStorage) Healtcheck() bool {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	err := s.db.PingContext(ctx)

	if err != nil {
		log.Printf("error pinging the database: %v", err)
	}
	return err == nil
}

// CreateTable - создает таблицу в базе данных
func (s *DBStorage) CreateTable(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, createTable)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%v", err)
	}

	return nil
}

// Put - добавляет запись в базу данных
func (s *DBStorage) Put(ctx context.Context, url string, uuid int) (string, error) {
	alias := generatestring.NewRandomString(8)

	// Проверяем, является ли строка JSON
	var jsonData map[string]string
	if err := json.Unmarshal([]byte(url), &jsonData); err == nil {
		// Если декодирование прошло успешно, используем значение "url"
		url = jsonData["url"]
	}

	shortURL := DBStorage{
		UUID:      uuid,
		URL:       url,
		Alias:     alias,
		CreatedAt: time.Now(),
	}

	_, err := s.db.ExecContext(ctx, insert, shortURL.UUID, shortURL.URL, shortURL.Alias, shortURL.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			// В случае конфликта выполняем дополнительный запрос для получения алиаса
			var existingAlias string
			err := s.db.QueryRowContext(ctx, getConflict, url).Scan(&existingAlias)
			if err != nil {
				if err == sql.ErrNoRows {
					// Обработка ситуации, когда не найдено совпадение по URL
					return alias, ErrConflict
				}
				s.log.Error("DBStorage.QueryRowContext failed", zap.Error(err))
				return alias, fmt.Errorf("failed to query existing alias: %w", err)
			}
			return existingAlias, ErrConflict
		}

		// Логирование текста ошибки для анализа
		s.log.Error("DBStorage.Put failed", zap.Error(err))
		return alias, fmt.Errorf("failed to insert short URL into database: %w", err)
	}

	// Логирование успешной вставки
	s.log.Info("DBStorage.Put", zap.String("url", url), zap.String("alias", alias))
	return alias, nil
}

// Get - получает URL по алиасу
func (s *DBStorage) Get(ctx context.Context, alias string) (string, bool, error) {
	var url DBStorage
	err := s.db.QueryRowContext(ctx, get, alias).Scan(&url.UUID, &url.URL, &url.Alias, &url.CreatedAt, &url.DeletedFlag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Это ожидаемая ошибка, когда нет строк, соответствующих запросу.
			return "", false, fmt.Errorf("no URL found for alias %s", alias)
		}
		// Обработка других ошибок базы данных
		return "", false, err
	}
	return url.URL, url.DeletedFlag, nil
}

// GetAll - получает все URL конкретного пользователя
func (s *DBStorage) GetAll(ctx context.Context, id int, host string) ([]byte, error) {
	var urls []DBStorage
	db := bun.NewDB(s.db, pgdialect.New())

	rows, err := db.NewSelect().
		TableExpr("short_urls").
		Column("url", "alias").
		Where("uuid = ?", id).
		Rows(ctx)
	if err != nil {
		s.log.Error("Error getting data: ", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		s.log.Error("Error iterating over rows: ", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	for rows.Next() {
		var url DBStorage
		if err := rows.Scan(&url.URL, &url.Alias); err != nil {
			s.log.Error("Error scanning data: ", zap.Error(err))
			return nil, err
		}
		url.Alias = fmt.Sprintf("%s/%s", host, url.Alias)
		urls = append(urls, url)
	}

	allURL, err := json.Marshal(urls)
	if err != nil {
		s.log.Error("Can't marshal URLs: ", zap.Error(err))
		return nil, err
	}

	return allURL, nil
}

// Del удаляет URL для пользователей с определенным ID.
func (s *DBStorage) Del(ctx context.Context, id int, aliases []string) error {
	inputCh := delGenerator(ctx, aliases)
	if err := s.updateStatusDel(ctx, id, inputCh); err != nil {
		return fmt.Errorf("failed to delete URLs: %w", err)
	}
	return nil
}

// Генератор канала для сбора alias'ов.
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

// updateStatusDel - выполняет пакетное обновление статуса удаления для заданных alias'ов и пользователя.
func (s *DBStorage) updateStatusDel(ctx context.Context, id int, inputCh chan string) error {
	var deleteURLs []string

	// Канал для сигнализации об окончании работы каждой горутины.
	aliasCollectionDoneCh := make(chan struct{})

	// Запуск горутины для сбора alias'ов.
	go func() {
		defer close(aliasCollectionDoneCh)
		for alias := range inputCh {
			deleteURLs = append(deleteURLs, alias)
		}
	}()

	// Ожидание окончания работы горутины по сбору alias'ов.
	<-aliasCollectionDoneCh

	db := bun.NewDB(s.db, pgdialect.New())

	// Пакетное обновление.
	_, err := db.NewUpdate().
		TableExpr("short_urls").
		Set("del = ?", true).
		Where("alias IN (?)", bun.In(deleteURLs)).
		Where("uuid = ?", id).
		Exec(ctx)

	if err != nil {
		s.log.Error("Can't exec update request: ", zap.Error(err))
		return fmt.Errorf("failed to update URLs: %w", err)
	}

	return nil
}

// GetAll - получает все URL конкретного пользователя (вариант 2)
// func (s *DBStorage) GetAll(ctx context.Context, id int, host string) ([]byte, error) {
// 	var urls []ShortURL
// 	allID, err := s.db.QueryContext(ctx, getAll, id)
// 	if err != nil {
// 		s.log.Error("Error getting batch data: ", zap.Error(err))
// 		return nil, err
// 	}
// 	defer func() {
// 		_ = allID.Close()
// 		_ = allID.Err()
// 	}()
//
// 	for allID.Next() {
// 		var url ShortURL
// 		if err := allID.Scan(&url.URL, &url.Alias); err != nil {
// 			s.log.Error("Error scanning data: ", zap.Error(err))
// 			return nil, err
// 		}
// 		urls = append(urls, ShortURL{
// 			URL:   url.URL,
// 			Alias: fmt.Sprintf("%s/%s", host, url.Alias),
// 		})
// 	}
// 	allURL, err := json.Marshal(urls)
// 	if err != nil {
// 		s.log.Error("Can't marshal URLs: ", zap.Error(err))
// 		return nil, err
// 	}
// 	return allURL, nil
// }
