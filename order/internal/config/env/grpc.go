package env

import (
	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	InventoryGrpcPort string `env:"INVENTORY_GRPC_PORT,required"`
	PaymentGrpcPort   string `env:"PAYMENT_GRPC_PORT,required"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

func NewGrpcConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &grpcConfig{raw: raw}, nil
}

func (cfg *grpcConfig) InventoryGrpcPort() string {
	return cfg.raw.InventoryGrpcPort
}

func (cfg *grpcConfig) PaymentGrpcPort() string {
	return cfg.raw.PaymentGrpcPort
}
