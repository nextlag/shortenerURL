// Package usecase provides use cases for managing short URLs, including
// interactions with data storage through a repository interface.
package usecase

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
)

// Repository represents the interface for data storage.
//
//go:generate mockgen -destination=mocks.go -package=usecase github.com/nextlag/shortenerURL/internal/usecase Repository
type Repository interface {
	Get(ctx context.Context, alias string) (string, bool, error)
	GetAll(ctx context.Context, userID int, url string) ([]byte, error)
	Put(ctx context.Context, url string, alias string, uuid int) (string, error)
	Del(ctx context.Context, id int, aliases []string) error
	Healthcheck() (bool, error)
}

// UseCase provides the use cases for interacting with the repository.
type UseCase struct {
	repo Repository        // interface for the repository
	log  *zap.Logger       // logger instance
	cfg  config.HTTPServer // configuration for the HTTP server
	DB   *sql.DB           // database connection
}

// New creates a new instance of UseCase.
func New(r Repository, l *zap.Logger, cfg config.HTTPServer) *UseCase {
	var db *sql.DB
	return &UseCase{repo: r, log: l, cfg: cfg, DB: db}
}

// DoGet retrieves a URL by its alias.
func (uc *UseCase) DoGet(ctx context.Context, alias string) (string, bool, error) {
	return uc.repo.Get(ctx, alias)
}

// DoGetAll retrieves all URLs for a specific user.
func (uc *UseCase) DoGetAll(ctx context.Context, userID int, url string) ([]byte, error) {
	return uc.repo.GetAll(ctx, userID, url)
}

// DoPut saves a URL with a generated alias.
func (uc *UseCase) DoPut(ctx context.Context, url string, alias string, uuid int) (string, error) {
	return uc.repo.Put(ctx, url, alias, uuid)
}

// DoDel deletes URLs for a user with the specified ID.
func (uc *UseCase) DoDel(ctx context.Context, id int, aliases []string) error {
	return uc.repo.Del(ctx, id, aliases)
}

// DoHealthcheck checks the health of the repository.
func (uc *UseCase) DoHealthcheck() (bool, error) {
	return uc.repo.Healthcheck()
}
