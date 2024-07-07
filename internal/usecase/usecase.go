// Package usecase provides use cases for managing short URLs, including
// interactions with data storage through a repository interface.
package usecase

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
)

// Repository represents the interface for data storage.
//
//go:generate mockgen -destination=mocks.go -package=usecase github.com/nextlag/shortenerURL/internal/usecase Repository
type Repository interface {
	Get(ctx context.Context, alias string) (string, bool, error)
	GetAll(ctx context.Context, userID int, url string) ([]byte, error)
	Put(ctx context.Context, url string, alias string, uuid int) (string, error)
	Del(id int, aliases []string) error
	Healthcheck() (bool, error)
}

// UseCase provides the use cases for interacting with the repository.
type UseCase struct {
	repo Repository            // interface for the repository
	log  *zap.Logger           // logger instance
	cfg  *configuration.Config // configuration for the HTTP server
	DB   *sql.DB               // database connection
}

// New creates a new instance of UseCase.
func New(r Repository, l *zap.Logger, cfg *configuration.Config) *UseCase {
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
func (uc *UseCase) DoDel(id int, aliases []string) {
	err := uc.repo.Del(id, aliases)
	if err != nil {
		uc.log.Error("Error deleting user URL", zap.Error(err))
	}
}

// DoHealthcheck checks the health of the repository.
func (uc *UseCase) DoHealthcheck() (bool, error) {
	return uc.repo.Healthcheck()
}
