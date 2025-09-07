package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHTTPEnvConfig struct {
	ShutdownTimeout   string `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
	HttpHost          string `env:"HTTP_HOST,required"`
	HttpPort          string `env:"HTTP_PORT,required"`
	ReadHeaderTimeout string `env:"HTTP_READ_HEADER_TIMEOUT,required"`
}

type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderHTTPConfig{raw: raw}, nil
}

func (cfg *orderHTTPConfig) ShutdownTimeout() time.Duration {
	duration, err := time.ParseDuration(cfg.raw.ShutdownTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return duration
}

func (cfg *orderHTTPConfig) HttpHost() string {
	return cfg.raw.HttpHost
}

func (cfg *orderHTTPConfig) HttpPort() string {
	return cfg.raw.HttpPort
}

func (cfg *orderHTTPConfig) ReadHeaderTimeout() time.Duration {
	duration, err := time.ParseDuration(cfg.raw.ReadHeaderTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return duration
}
