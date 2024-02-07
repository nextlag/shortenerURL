package usecase

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
)

// Repository представляет интерфейс для хранилища данных
//
//go:generate mockgen -destination=mocks.go -package=usecase github.com/nextlag/shortenerURL/internal/usecase Repository
type Repository interface {
	Get(ctx context.Context, alias string) (string, bool, error)
	GetAll(ctx context.Context, userID int, url string) ([]byte, error)
	Put(ctx context.Context, url string, uuid int) (string, error)
	Del(ctx context.Context, id int, aliases []string) error
	Healthcheck() (bool, error)
}

type UseCase struct {
	repo Repository // interface Repository
	log  *zap.Logger
	cfg  config.HTTPServer
	DB   *sql.DB
}

func New(r Repository, l *zap.Logger, cfg config.HTTPServer) *UseCase {
	var db *sql.DB
	return &UseCase{repo: r, log: l, cfg: cfg, DB: db}
}

func (uc *UseCase) DoGet(ctx context.Context, alias string) (string, bool, error) {
	return uc.repo.Get(ctx, alias)
}

func (uc *UseCase) DoGetAll(ctx context.Context, userID int, url string) ([]byte, error) {
	return uc.repo.GetAll(ctx, userID, url)
}

func (uc *UseCase) DoPut(ctx context.Context, url string, uuid int) (string, error) {
	return uc.repo.Put(ctx, url, uuid)
}

func (uc *UseCase) DoDel(ctx context.Context, id int, aliases []string) error {
	return uc.repo.Del(ctx, id, aliases)
}

func (uc *UseCase) DoHealthcheck() (bool, error) {
	return uc.repo.Healthcheck()
}
