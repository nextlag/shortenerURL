package psql

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

	"github.com/nextlag/shortenerURL/internal/configuration"
	"github.com/nextlag/shortenerURL/internal/entity"
	"github.com/nextlag/shortenerURL/internal/usecase/repository/models"
	"github.com/nextlag/shortenerURL/pkg/tools/generatestring"
)

type Repo struct {
	DB  *sql.DB
	cfg *configuration.Config
	log *zap.Logger
}

// New creates a new Repo instance.
func New(cfg *configuration.Config, log *zap.Logger) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), createTablesTimeout)
	defer cancel()

	DB, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		log.Error("error when opening a connection to the database", zap.Error(err))
		return nil, fmt.Errorf("DB connection error: %w", err)
	}

	if err := DB.PingContext(ctx); err != nil {
		log.Error("error when checking database connection", zap.Error(err))
		return nil, fmt.Errorf("DB ping error: %w", err)
	}

	r := &Repo{
		DB:  DB,
		cfg: cfg,
		log: log,
	}

	if err = r.CreateTable(ctx); err != nil {
		return nil, fmt.Errorf("create table error: %w", err)
	}

	return r, nil
}

var ErrConflict = errors.New("data conflict in DBStorage")

const (
	pingTimeout         = time.Second * 3
	createTablesTimeout = time.Second * 5
	createTable         = `CREATE TABLE IF NOT EXISTS short_urls (
		uuid INT,
		url VARCHAR NOT NULL,
		alias VARCHAR(255) NOT NULL,
		created_at TIMESTAMP,
		del BOOLEAN,
		PRIMARY KEY (uuid, alias),
		UNIQUE (uuid, url)
	);`
	insert       = `INSERT INTO short_urls (uuid, url, alias, created_at, del) VALUES ($1, $2, $3, $4, false);`
	get          = `SELECT uuid, url, alias, created_at, del FROM short_urls WHERE alias = $1;`
	getConflict  = `SELECT alias FROM short_urls WHERE url = $1;`
	getUrlsStats = `SELECT COUNT(*) as urlsCount FROM short_urls;`
	getUserStats = `SELECT COUNT(DISTINCT uuid) as uniqueUsers FROM short_urls;`
)

// Stop closes the connection to the database.
func (r *Repo) Stop() error {
	return r.DB.Close()
}

// Healthcheck checks the connection to the database.
func (r *Repo) Healthcheck() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := r.DB.PingContext(ctx); err != nil {
		r.log.Error("Ошибка подключения к базе данных", zap.Error(err))
		return false, err
	}
	return true, nil
}

// CreateTable creates the short_urls table in the database.
func (r *Repo) CreateTable(ctx context.Context) error {
	_, err := r.DB.ExecContext(ctx, createTable)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%v", err)
	}
	return nil
}

// Put adds a new short URL to the database.
func (r *Repo) Put(ctx context.Context, url string, alias string, userID int) (string, error) {
	if alias == "" {
		alias = generatestring.NewRandomString(8)
	}

	var jsonData map[string]string
	if err := json.Unmarshal([]byte(url), &jsonData); err == nil {
		url = jsonData["url"]
	}

	shortURL := entity.URL{
		UUID:      userID,
		URL:       url,
		Alias:     alias,
		CreatedAt: time.Now(),
		IsDeleted: false,
	}

	_, err := r.DB.ExecContext(ctx, insert, shortURL.UUID, shortURL.URL, shortURL.Alias, shortURL.IsDeleted, shortURL.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			var existingAlias string
			err = r.DB.QueryRowContext(ctx, getConflict, url).Scan(&existingAlias)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
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

// Get retrieves a URL by its alias.
func (r *Repo) Get(ctx context.Context, alias string) (*entity.URL, error) {
	var url entity.URL
	err := r.DB.QueryRowContext(ctx, get, alias).Scan(&url.UUID, &url.URL, &url.Alias, &url.CreatedAt, &url.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no URL found for alias %s", alias)
		}
		return nil, err
	}
	return &url, nil
}

// GetAll retrieves all URLs for a specific user.
func (r *Repo) GetAll(ctx context.Context, userID int, host string) ([]*entity.URL, error) {
	var urls []*entity.URL
	DB := bun.NewDB(r.DB, pgdialect.New())

	rows, err := DB.NewSelect().
		TableExpr("short_urls").
		Column("url", "alias", "del", "created_at").
		Where("uuid = ?", userID).
		Rows(ctx)
	if err != nil {
		r.log.Error("Error getting data: ", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url entity.URL
		if err = rows.Scan(&url.URL, &url.Alias, &url.IsDeleted, &url.CreatedAt); err != nil {
			r.log.Error("Error scanning data: ", zap.Error(err))
			return nil, err
		}
		url.Alias = fmt.Sprintf("%s/%s", host, url.Alias)
		urls = append(urls, &url)
	}

	if rows.Err() != nil {
		r.log.Error("Error iterating over rows: ", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return urls, nil
}

// Del removes URLs for a user with a specific ID.
func (r *Repo) Del(ctx context.Context, userID int, aliases []string) error {
	DB := bun.NewDB(r.DB, pgdialect.New())

	_, err := DB.NewUpdate().
		TableExpr("short_urls").
		Set("del = ?", true).
		Where("alias IN (?)", bun.In(aliases)).
		Where("uuid = ?", userID).
		Exec(ctx)

	if err != nil {
		r.log.Error("Can't exec update request: ", zap.Error(err))
		return fmt.Errorf("failed to update URLs: %w", err)
	}

	return nil
}

// GetStats retrieves statistics on users and URLs.
func (r *Repo) GetStats(ctx context.Context) ([]byte, error) {
	urlsStatRaw := r.DB.QueryRowContext(ctx, getUrlsStats)
	userStatRaw := r.DB.QueryRowContext(ctx, getUserStats)

	var urlsStat, usersStat int

	if err := urlsStatRaw.Scan(&urlsStat); err != nil {
		return nil, fmt.Errorf("error scanning urlsStat: %w", err)
	}

	if err := userStatRaw.Scan(&usersStat); err != nil {
		return nil, fmt.Errorf("error scanning usersStat: %w", err)
	}

	resultStats, err := json.Marshal(models.Stats{
		URLs:  urlsStat,
		Users: usersStat,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling Stats: %w", err)
	}

	return resultStats, nil
}
