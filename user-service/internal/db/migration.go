package db

import (
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source driver
)

func RunMigration(dbUrl string) {

	_, b, _, _ := runtime.Caller(0) // get current file path
	basepath := filepath.Dir(b)     // directory of migration.go
	migrationsPath := "file://" + filepath.Join(basepath, "../../migrations")

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
