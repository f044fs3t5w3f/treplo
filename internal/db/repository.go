package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	db *sql.DB
}

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

func (r *repository) AddAudio(ctx context.Context, chatID int64, fileID string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO files (chat_id, file_id)
VALUES ($1, $2)`, chatID, fileID)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}
	return nil
}

func migrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
