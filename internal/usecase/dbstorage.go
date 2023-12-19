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
	CreateTable = `CREATE TABLE IF NOT EXISTS short_urls (
    uuid INT,
    url VARCHAR NOT NULL,
    alias VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    del BOOLEAN,
    PRIMARY KEY (uuid, alias),
    UNIQUE (uuid, url)
);`

	Insert      = `INSERT INTO short_urls (uuid, url, alias, created_at, del) VALUES ($1, $2, $3, $4, false);`
	Get         = `SELECT uuid, url, alias, created_at, del FROM short_urls WHERE alias = $1;`
	GetConflict = `SELECT alias FROM short_urls WHERE url = $1;`
)

const createTablesTimeout = time.Second * 5

// NewDB - создает новый экземпляр DBStorage
func NewDB(cfg string, log *zap.Logger) (*UseCase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), createTablesTimeout)
	defer cancel()
	// Создание подключения к базе данных с использованием контекста
	DB, err := sql.Open("pgx", cfg)
	if err != nil {
		log.Error("error when opening a connection to the database", zap.Error(err))
		return nil, fmt.Errorf("DB connection err=%w", err)
	}
	// Проверка подключения к базе данных с использованием контекста
	if err := DB.PingContext(ctx); err != nil {
		log.Error("error when checking database connection", zap.Error(err))
		return nil, fmt.Errorf("DB ping error: %w", err)
	}
	storage := &UseCase{
		DB: DB,
	}
	// Создание таблицы с использованием контекста
	if err := storage.CreateTable(ctx); err != nil {
		return nil, fmt.Errorf("create table error: %w", err)
	}
	return storage, nil
}

// Время ожидания пинга для проверки подключения к базе данных
const pingTimeout = time.Second * 3

// Stop - закрывает соединение с базой данных
func (uc *UseCase) Stop() error {
	uc.DB.Close()
	return nil
}

// Healthcheck - проверяет подключение к базе данных
func (uc *UseCase) Healthcheck() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := uc.DB.PingContext(ctx); err != nil {
		// Обработка ошибки подключения, например, вывод в лог или возврат false
		uc.log.Error("Ошибка подключения к базе данных", zap.Error(err))
		return false, err
	}
	// Подключение успешно
	return true, nil
}

// CreateTable - создает таблицу в базе данных
func (uc *UseCase) CreateTable(ctx context.Context) error {
	_, err := uc.DB.ExecContext(ctx, CreateTable)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%v", err)
	}

	return nil
}

// Put - добавляет запись в базу данных
func (uc *UseCase) Put(ctx context.Context, url string, uuid int) (string, error) {
	alias := generatestring.NewRandomString(8)

	// Проверяем, является ли строка JSON
	var jsonData map[string]string
	if err := json.Unmarshal([]byte(url), &jsonData); err == nil {
		// Если декодирование прошло успешно, используем значение "url"
		url = jsonData["url"]
	}

	shortURL := entity.DBStorage{
		UUID:      uuid,
		URL:       url,
		Alias:     alias,
		CreatedAt: time.Now(),
	}

	_, err := uc.DB.ExecContext(ctx, Insert, shortURL.UUID, shortURL.URL, shortURL.Alias, shortURL.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			// В случае конфликта выполняем дополнительный запрос для получения алиаса
			var existingAlias string
			err := uc.DB.QueryRowContext(ctx, GetConflict, url).Scan(&existingAlias)
			if err != nil {
				if err == sql.ErrNoRows {
					// Обработка ситуации, когда не найдено совпадение по URL
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

// Get - получает URL по алиасу
func (uc *UseCase) Get(ctx context.Context, alias string) (string, bool, error) {
	var url entity.DBStorage
	err := uc.DB.QueryRowContext(ctx, Get, alias).Scan(&url.UUID, &url.URL, &url.Alias, &url.CreatedAt, &url.DeletedFlag)
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

	if rows.Err() != nil {
		uc.log.Error("Error iterating over rows: ", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	for rows.Next() {
		var url entity.DBStorage
		if err := rows.Scan(&url.URL, &url.Alias); err != nil {
			uc.log.Error("Error scanning data: ", zap.Error(err))
			return nil, err
		}
		url.Alias = fmt.Sprintf("%s/%s", host, url.Alias)
		urls = append(urls, url)
	}

	allURL, err := json.Marshal(urls)
	if err != nil {
		uc.log.Error("Can't marshal URLs: ", zap.Error(err))
		return nil, err
	}

	return allURL, nil
}

// Del удаляет URL для пользователей с определенным ID.
func (uc *UseCase) Del(ctx context.Context, id int, aliases []string) error {
	inputCh := delGenerator(ctx, aliases)
	if err := uc.updateStatusDel(ctx, id, inputCh); err != nil {
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
func (uc *UseCase) updateStatusDel(ctx context.Context, id int, inputCh chan string) error {
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

	DB := bun.NewDB(uc.DB, pgdialect.New())

	// Пакетное обновление.
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
