package db

import (
	"fmt"

	"github.com/auroban/gochat-be/auth-service/internal/config"
	"github.com/jmoiron/sqlx"
)

func Connect(cfg config.Config) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
		cfg.Database.Schema,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
