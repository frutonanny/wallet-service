package postgres

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const pathMigration = "migrations"

// MustConnect - устанавливает соединение с базой.
func MustConnect(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(fmt.Errorf("opening pgx driver: %w", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("connect db: %w", err))
	}

	return db
}

// MustMigrate - применяет миграции из переданной директории
func MustMigrate(db *sql.DB) {
	if err := goose.Up(db, pathMigration); err != nil {
		panic(fmt.Errorf("apply migrations: %w", err))
	}
}
