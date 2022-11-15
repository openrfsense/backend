package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/openrfsense/common/logging"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
)

var log = logging.New().
	WithPrefix("storage").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

const (
	maxPoolSize        = 100
	connectionAttempts = 10
	connectionTimeout  = time.Second
)

var pg *Postgres

// Type Postgres contains a database connection pool and a query builder.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	squirrel.StatementBuilderType
	*pgxpool.Pool
}

// Initalizes a connection to a PostgreSQL database.
func Init(config *koanf.Koanf) error {
	pg = &Postgres{
		maxPoolSize:  maxPoolSize,
		connAttempts: connectionAttempts,
		connTimeout:  connectionTimeout,
	}

	pg.StatementBuilderType = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(generateConnString(config))
	if err != nil {
		return err
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	ticker := time.NewTicker(pg.connTimeout)
	defer ticker.Stop()
	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Debugf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		// time.Sleep(pg.connTimeout)
		<-ticker.C
		pg.connAttempts--
	}

	if err != nil {
		return err
	}

	return autoMigrate(config)
}

// Closes the database connection.
func Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}

// Returns a reference to the internal database connection and query builder.
func Instance() *Postgres {
	return pg
}

// Generates a valid connection string from the configuration
func generateConnString(config *koanf.Koanf) string {
	slice := []string{}

	slice = append(slice, "host="+config.MustString("postgres.host"))
	slice = append(slice, fmt.Sprintf("port=%d", config.MustInt("postgres.port")))
	slice = append(slice, "user="+config.MustString("postgres.username"))
	slice = append(slice, "password="+config.MustString("postgres.password"))

	slice = append(slice, config.Strings("postgres.params")...)

	return strings.Join(slice, " ")
}
