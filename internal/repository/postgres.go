package repository

import (
	"database/sql"
	"fmt"

	"current-account-service/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres(cfg config.Config) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
