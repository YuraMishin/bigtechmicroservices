package env

import (
	"strconv"

	"github.com/caarlos0/env/v11"
)

type paymentGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
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
		// Возвращаем значение по умолчанию в случае ошибки
		return 0
	}
	return port
}
