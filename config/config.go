package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

type Backend struct {
	Metrics bool              `yaml:"metrics"`
	Port    int               `yaml:"port"`
	Users   map[string]string `yaml:"users"`
}

type Postgres struct {
	Host         string   `yaml:"host"`
	Port         int      `yaml:"port"`
	DatabaseName string   `yaml:"dbname"`
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password"`
	Params       []string `yaml:"params"`
}

type NATS struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Secret   string `yaml:"secret"`
}

type BackendConfig struct {
	Backend  `yaml:"backend"`
	Postgres `yaml:"postgres"`
	NATS     `yaml:"nats"`
}

var defaultConfig = BackendConfig{
	Backend: Backend{
		Metrics: true,
		Port:    8080,
	},
	Postgres: Postgres{
		Host:         "localhost",
		Port:         5432,
		DatabaseName: "postgres",
		Username:     "postgres",
		Password:     "postgres",
	},
	NATS: NATS{
		Protocol: "tcp",
		Port:     4222,
	},
}

var konf *koanf.Koanf

// Loads a YAML configuration file from the given path and overrides
// it with environment variables. If the file cannot be loaded or
// parsed as YAML, an error is returned. Requires a default config of any kind,
// will try to serialize the configuration to outConfig if present (needs to
// be a pointer to a struct).
func Load(path string) (*koanf.Koanf, error) {
	konf = koanf.New(".")

	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("configuration file path cannot be empty")
	}

	_ = konf.Load(structs.Provider(defaultConfig, "yaml"), nil)

	if err := konf.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading configuration file: %v (%T)", err, err)
	}

	_ = konf.Load(env.Provider("ORFS_", ".", formatEnv), nil)

	return konf, nil
}

// Formats environment variables: ORFS_SECTION_SUBSECTION_KEY becomes
// (as a path) section.subsection.key
func formatEnv(s string) string {
	rawPath := strings.ToLower(strings.TrimPrefix(s, "ORFS_"))
	return strings.Replace(rawPath, "_", ".", -1)
}
