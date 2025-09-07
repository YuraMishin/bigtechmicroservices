package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderGRPCConfig interface {
	ShutdownTimeout() time.Duration
	InventoryGrpcHost() string
	InventoryGrpcPort() string
	PaymentGrpcHost() string
	PaymentGrpcPort() string
}

type OrderHTTPConfig interface {
	ShutdownTimeout() time.Duration
	HttpHost() string
	HttpPort() string
	ReadHeaderTimeout() time.Duration
}

type OrderPostgresConfig interface {
	PostgresHost() string
	PostgresPort() string
	PostgresUser() string
	PostgresPassword() string
	PostgresDB() string
	PostgresSSLMode() string
	MigrationDirectory() string
}
