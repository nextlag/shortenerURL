package dbstorage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
)

type DBStorage struct {
	db        *sql.DB
	log       *zap.Logger
	UUID      int       `json:"user_id,omitempty" `
	URL       string    `json:"original_url,omitempty"`
	Alias     string    `json:"short_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Время ожидания пинга для проверки подключения к базе данных
const pingTimeout = time.Second * 3

// Время ожидания создания таблицы
const createTablesTimeout = time.Second * 5

// ErrConflict - ошибка конфликта данных
var ErrConflict = errors.New("data conflict in DBStorage")

// New - создает новый экземпляр DBStorage
func New(cfg string, log *zap.Logger) (*DBStorage, error) {
	db, err := sql.Open("pgx", cfg)
	if err != nil {
		return nil, fmt.Errorf("db connection err=%w", err)
	}
	storage := &DBStorage{
		db:  db,
		log: log,
	}
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
func (s *DBStorage) Get(ctx context.Context, alias string) (string, error) {
	var url DBStorage
	err := s.db.QueryRowContext(ctx, get, alias).Scan(&url.UUID, &url.URL, &url.Alias, &url.CreatedAt)
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

// Del - Удаление URL для пользователей с определенным ID
func (s *DBStorage) DeleteURLs(_ context.Context, id int, shortURLs []string) error {
	context := context.Background()
	inputCh := deleteURLsGenerator(context, shortURLs)
	s.bulkDeleteStatusUpdate(id, inputCh)
	return nil
}

func deleteURLsGenerator(ctx context.Context, URLs []string) chan string {
	URLsCh := make(chan string)
	go func() {
		defer close(URLsCh)
		for _, data := range URLs {
			select {
			case <-ctx.Done():
				return
			case URLsCh <- data:
			}
		}
	}()
	return URLsCh
}

func (s *DBStorage) bulkDeleteStatusUpdate(id int, inputChs ...chan string) {
	var wg sync.WaitGroup

	deleteUpdate := func(c chan string) {
		var linksToDelete []string
		for shortenLink := range c {
			linksToDelete = append(linksToDelete, shortenLink)
		}
		db := bun.NewDB(s.db, pgdialect.New())

		_, err := db.NewUpdate().
			TableExpr("shorten_URLs").
			Set("deleted = ?", "true").
			Where("short_link IN (?)", bun.In(linksToDelete)).
			WhereGroup(" AND ", func(uq *bun.UpdateQuery) *bun.UpdateQuery {
				return uq.Where("user_id = ?", id)
			}).
			Exec(context.Background())
		if err != nil {
			s.log.Error("Can't exec update request: ", zap.Error(err))
		}
		wg.Done()
	}

	wg.Add(len(inputChs))

	for _, c := range inputChs {
		go deleteUpdate(c)
	}

	go func() {
		wg.Wait()
	}()
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
