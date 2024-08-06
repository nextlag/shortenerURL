// Package usecase provides use cases for managing short URLs, including
// interactions with data storage through a repository interface.
package usecase

import (
	"context"
	"fmt"

	"github.com/nextlag/shortenerURL/internal/entity"
	"github.com/nextlag/shortenerURL/internal/usecase/repository"
)

// UseCase provides the use cases for interacting with the repository.
type UseCase struct {
	repo repository.Repository // interface for the repository
}

// New creates a new instance of UseCase.
func New(r repository.Repository) *UseCase {
	return &UseCase{repo: r}
}

// DoGet retrieves a URL by its alias.
func (uc *UseCase) DoGet(ctx context.Context, alias string) (*entity.URL, error) {
	return uc.repo.Get(ctx, alias)
}

// DoGetAll retrieves all URLs for a specific user.
func (uc *UseCase) DoGetAll(ctx context.Context, userID int, url string) ([]*entity.URL, error) {
	return uc.repo.GetAll(ctx, userID, url)
}

// DoPut saves a URL with a generated alias.
func (uc *UseCase) DoPut(ctx context.Context, url string, alias string, uuid int) (string, error) {
	return uc.repo.Put(ctx, url, alias, uuid)
}

// DoDel deletes URLs for a user with the specified ID.
func (uc *UseCase) DoDel(ctx context.Context, id int, aliases []string) {
	err := uc.repo.Del(ctx, id, aliases)
	if err != nil {
		_ = fmt.Errorf("error deleting user URL: %w", err)
	}
}

// DoHealthcheck checks the health of the repository.
func (uc *UseCase) DoHealthcheck() (bool, error) {
	return uc.repo.Healthcheck()
}

// DoGetStats requests to the database to obtain statistics
func (uc *UseCase) DoGetStats(ctx context.Context) ([]byte, error) {
	return uc.repo.GetStats(ctx)
}
