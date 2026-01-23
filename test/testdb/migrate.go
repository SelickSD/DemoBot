package testdb

import (
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {
	goose.SetDialect("postgres")

	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// поднимаемся к корню проекта
	projectRoot := filepath.Clean(filepath.Join(wd, "../../.."))

	migrationsDir := filepath.Join(projectRoot, "build", "app", "migrations")

	return goose.Up(db, migrationsDir)
}
