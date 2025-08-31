package env

import (
	"strconv"

	"github.com/caarlos0/env/v11"
)

type inventoryGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
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
		// Возвращаем значение по умолчанию в случае ошибки
		return 0
	}
	return port
}
