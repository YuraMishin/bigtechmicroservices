package env

import (
	"github.com/caarlos0/env/v11"
)

type orderPostgresEnvConfig struct {
	PostgresHost       string `env:"POSTGRES_HOST,required"`
	PostgresPort       string `env:"POSTGRES_PORT,required"`
	PostgresUser       string `env:"POSTGRES_USER,required"`
	PostgresPassword   string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB         string `env:"POSTGRES_DB,required"`
	PostgresSSLMode    string `env:"POSTGRES_SSL_MODE,required"`
	MigrationDirectory string `env:"MIGRATION_DIRECTORY,required"`
}

type orderPostgresConfig struct {
	raw orderPostgresEnvConfig
}

func NewOrderPostgresConfig() (*orderPostgresConfig, error) {
	var raw orderPostgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPostgresConfig{raw: raw}, nil
}

func (cfg *orderPostgresConfig) PostgresHost() string {
	return cfg.raw.PostgresHost
}

func (cfg *orderPostgresConfig) PostgresPort() string {
	return cfg.raw.PostgresPort
}

func (cfg *orderPostgresConfig) PostgresUser() string {
	return cfg.raw.PostgresUser
}

func (cfg *orderPostgresConfig) PostgresPassword() string {
	return cfg.raw.PostgresPassword
}

func (cfg *orderPostgresConfig) PostgresDB() string {
	return cfg.raw.PostgresDB
}

func (cfg *orderPostgresConfig) PostgresSSLMode() string {
	return cfg.raw.PostgresSSLMode
}

func (cfg *orderPostgresConfig) MigrationDirectory() string {
	return cfg.raw.MigrationDirectory
}
