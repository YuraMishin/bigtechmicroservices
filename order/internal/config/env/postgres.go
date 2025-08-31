package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	PostgresHost       string `env:"POSTGRES_HOST,required"`
	PostgresPort       string `env:"POSTGRES_PORT,required"`
	PostgresDB         string `env:"POSTGRES_DB,required"`
	PostgresUser       string `env:"POSTGRES_USER,required"`
	PostgresPassword   string `env:"POSTGRES_PASSWORD,required"`
	PostgresSSLMode    string `env:"POSTGRES_SSL_MODE,required"`
	MigrationDirectory string `env:"MIGRATION_DIRECTORY,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.raw.PostgresUser,
		cfg.raw.PostgresPassword,
		cfg.raw.PostgresHost,
		cfg.raw.PostgresPort,
		cfg.raw.PostgresDB,
		cfg.raw.PostgresSSLMode,
	)
}

func (cfg *postgresConfig) MigrationDirectory() string {
	return cfg.raw.MigrationDirectory
}
