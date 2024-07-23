package repository

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
	"github.com/nextlag/shortenerURL/internal/entity"
	"github.com/nextlag/shortenerURL/internal/usecase/repository/inmemory"
	"github.com/nextlag/shortenerURL/internal/usecase/repository/psql"
)

// Repository represents the interface for data storage.
//
//go:generate mockgen -destination=mocks.go -package=repository github.com/nextlag/shortenerURL/internal/usecase/repository Repository
type Repository interface {
	Get(ctx context.Context, alias string) (*entity.URL, error)
	GetAll(ctx context.Context, userID int, host string) ([]*entity.URL, error)
	Put(ctx context.Context, url string, alias string, userID int) (string, error)
	Del(ctx context.Context, userID int, aliases []string) error
	Healthcheck() (bool, error)
	GetStats(ctx context.Context) ([]byte, error)
}

const (
	postgres = "pg"
	inMemory = "mem"
)

// New repository
func New(cfg *configuration.Config, log *zap.Logger) (Repository, error) {
	switch cfg.StorageType {
	case postgres:
		if cfg.DSN != "" && cfg.FileStorage == "" {
			db, err := psql.New(cfg, log)
			if err != nil {
				log.Fatal("failed to connect to database", zap.Error(err))
			}
			return db, nil
		} else {
			log.Fatal("the configuration is incorrect: remove the FileStorage configuration parameter")
		}
	case inMemory:
		if cfg.FileStorage != "" && cfg.DSN == "" {
			db, err := inmemory.New(cfg, log)
			if err != nil {
				log.Fatal("failed in memory storage", zap.Error(err))
			}
			err = inmemory.Load(cfg.FileStorage, db)
			if err != nil {
				log.Fatal("failed to load data from file", zap.Error(err))
			}
			return db, nil
		} else {
			log.Fatal("the configuration is incorrect: remove the DSN configuration parameter")
		}
	default:
		return nil, errors.New("unknown storage type")
	}
	return nil, errors.New("configuration error")
}
