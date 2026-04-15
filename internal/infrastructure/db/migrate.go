package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		_, _ = m.Close()
		return fmt.Errorf("apply migrations: %w", err)
	}

	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("close migration source: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("close migration database: %w", dbErr)
	}

	return nil
}
