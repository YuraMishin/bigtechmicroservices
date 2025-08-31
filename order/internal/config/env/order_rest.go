package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type orderRESTEnvConfig struct {
	Port              string `env:"HTTP_PORT,required"`
	ReadHeaderTimeout string `env:"HTTP_READ_HEADER_TIMEOUT,required"`
	ShutdownTimeout   string `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
}

type orderRESTConfig struct {
	raw orderRESTEnvConfig
}

func NewOrderRESTConfig() (*orderRESTConfig, error) {
	var raw orderRESTEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderRESTConfig{raw: raw}, nil
}

func (cfg *orderRESTConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *orderRESTConfig) ReadHeaderTimeout() time.Duration {
	d, err := time.ParseDuration(cfg.raw.ReadHeaderTimeout)
	if err != nil {
		return 0
	}
	return d
}

func (cfg *orderRESTConfig) ShutdownTimeout() time.Duration {
	d, err := time.ParseDuration(cfg.raw.ShutdownTimeout)
	if err != nil {
		return 0
	}
	return d
}
