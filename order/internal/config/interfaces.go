package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderRESTConfig interface {
	Port() string
	ReadHeaderTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}

type GRPCConfig interface {
	InventoryGrpcPort() string
	PaymentGrpcPort() string
}
