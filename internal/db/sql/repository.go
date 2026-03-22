package sql

import (
	"database/sql"
	"fmt"

	"github.com/a-kuleshov/treplo/internal/db"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	db *sql.DB
}

var _ db.Repository = &repository{}

func NewRepository(databaseDSN string) (*repository, error) {
	db, err := sql.Open("pgx", databaseDSN)

	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	err = migrateDB(db)

	if err != nil {
		return nil, fmt.Errorf("migrateDB: %w", err)
	}

	return &repository{db: db}, nil
}
