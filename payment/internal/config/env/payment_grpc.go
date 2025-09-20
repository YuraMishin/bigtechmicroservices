package env

import (
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
)

type paymentGRPCEnvConfig struct {
	Host            string `env:"GRPC_HOST,required"`
	Port            string `env:"GRPC_PORT,required"`
	ShutdownTimeout string `env:"GRPC_SHUTDOWN_TIMEOUT,required"`
}

type paymentGRPCConfig struct {
	raw paymentGRPCEnvConfig
}

func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	var raw paymentGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentGRPCConfig{raw: raw}, nil
}

func (cfg *paymentGRPCConfig) Port() int {
	port, err := strconv.Atoi(cfg.raw.Port)
	if err != nil {
		return 8080
	}
	return port
}

func (cfg *paymentGRPCConfig) ShutdownTimeout() time.Duration {
	duration, err := time.ParseDuration(cfg.raw.ShutdownTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return duration
}
