package database

import (
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/knadh/koanf"
)

const (
	migrateAttempts = 5
	migrateTimeout  = time.Second
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

// Performs migrations as needed on the configured database.
func autoMigrate(config *koanf.Koanf) error {
	var (
		attempts = migrateAttempts
		err      error
		m        *migrate.Migrate
	)

	migrations, err := iofs.New(migrationsFs, "migrations")
	defer func() {
		err = migrations.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.MustString("postgres.username"),
		config.MustString("postgres.password"),
		config.MustString("postgres.host"),
		config.MustInt("postgres.port"),
		config.MustString("postgres.dbname"),
	)

	for attempts > 0 {
		m, err = migrate.NewWithSourceInstance("iofs", migrations, connString)
		if err == nil {
			break
		}

		log.Debugf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(migrateTimeout)
		attempts--
	}

	if err != nil {
		return err
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Warn("Migrate: no change")
		return nil
	}

	log.Debug("Migrated successfully")
	return nil
}
