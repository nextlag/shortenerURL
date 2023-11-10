package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const pingTimeout = time.Second * 10

type DBStorage struct {
	db *sql.DB
}

func New(dataSourceName string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("db conntection err=%w", err)
	}

	return &DBStorage{
		db: db,
	}, nil
}

func (s *DBStorage) Stop() error {
	s.db.Close()

	return nil
}

func (s *DBStorage) CheckConnection() bool {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	err := s.db.PingContext(ctx)

	if err != nil {
		log.Printf("print db err=%w", err)
	}

	return err == nil
}
