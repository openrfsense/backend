package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/openrfsense/common/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/knadh/koanf"
)

var log = logging.New().
	WithPrefix("storage").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

var pg *gorm.DB

// Initalizes a connection to a PostgreSQL database.
func Init(config *koanf.Koanf) error {
	var err error
	conn := postgres.Open(generateConnString(config))
	// Force timestamps to UTC
	pg, err = gorm.Open(conn, &gorm.Config{
		NowFunc: time.Now().UTC,
	})
	if err != nil {
		return err
	}

	err = pg.AutoMigrate(
		&Campaign{},
		&Sample{},
	)
	if err != nil {
		return err
	}

	log.Debug("Connected to Postgres DB")

	return nil
}

// Returns a reference to the internal databse connection.
func Instance() *gorm.DB {
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
