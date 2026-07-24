package db

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source driver
)

func RunMigration(dbUrl string) {
	// Default to ./migrations, can be overridden by MIGRATIONS_PATH env var
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	m, err := migrate.New(
		migrationsPath,
		dbUrl,
	)

	if err != nil {
		logger.Fatalf("❌ Failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatalf("❌ Migration failed: %v", err)
	}

	logger.Info("✅ Migrations applied successfully")

}
