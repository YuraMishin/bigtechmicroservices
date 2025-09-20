package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type orderGRPCEnvConfig struct {
	ShutdownTimeout   string `env:"GRPC_SHUTDOWN_TIMEOUT,required"`
	InventoryGrpcHost string `env:"INVENTORY_GRPC_HOST,required"`
	InventoryGrpcPort string `env:"INVENTORY_GRPC_PORT,required"`
	PaymentGrpcHost   string `env:"PAYMENT_GRPC_HOST,required"`
	PaymentGrpcPort   string `env:"PAYMENT_GRPC_PORT,required"`
}

type orderGRPCConfig struct {
	raw orderGRPCEnvConfig
}

func NewOrderGRPCConfig() (*orderGRPCConfig, error) {
	var raw orderGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderGRPCConfig{raw: raw}, nil
}

func (cfg *orderGRPCConfig) ShutdownTimeout() time.Duration {
	duration, err := time.ParseDuration(cfg.raw.ShutdownTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return duration
}

func (cfg *orderGRPCConfig) InventoryGrpcHost() string {
	return cfg.raw.InventoryGrpcHost
}

func (cfg *orderGRPCConfig) InventoryGrpcPort() string {
	return cfg.raw.InventoryGrpcPort
}

func (cfg *orderGRPCConfig) PaymentGrpcHost() string {
	return cfg.raw.PaymentGrpcHost
}

func (cfg *orderGRPCConfig) PaymentGrpcPort() string {
	return cfg.raw.PaymentGrpcPort
}
