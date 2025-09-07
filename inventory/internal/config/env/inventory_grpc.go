package env

import (
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
)

type inventoryGRPCEnvConfig struct {
	Host            string `env:"GRPC_HOST,required"`
	Port            string `env:"GRPC_PORT,required"`
	ShutdownTimeout string `env:"GRPC_SHUTDOWN_TIMEOUT,required"`
}

type inventoryGRPCConfig struct {
	raw inventoryGRPCEnvConfig
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{raw: raw}, nil
}

func (cfg *inventoryGRPCConfig) Port() int {
	port, err := strconv.Atoi(cfg.raw.Port)
	if err != nil {
		return 8080
	}
	return port
}

func (cfg *inventoryGRPCConfig) ShutdownTimeout() time.Duration {
	duration, err := time.ParseDuration(cfg.raw.ShutdownTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return duration
}
