package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrate(db *sql.DB) error {
	var (
		err error
		m   *migrate.Migrate
	)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate -  postgres.WithInstance: %v", err)
	}

	m, err = migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate - migrate.NewWithDatabaseInstance: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: up error: %s", err)
	}
	// defer m.Close()

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrate: no change")
	}

	log.Println("Migrate: up success")

	return nil
}
