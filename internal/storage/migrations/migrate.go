package migrations

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(connectionURL string) error {
	workdir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get work directory: %w", err)
	}

	rootPath, err := filepath.Abs(filepath.Join(workdir, "..", "..", "internal", "storage", "migrations"))
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	migrationURL := &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(rootPath),
	}

	m, err := migrate.New(migrationURL.String(), connectionURL)
	if err != nil {
		return fmt.Errorf("failed to create migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
