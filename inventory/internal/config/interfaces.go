package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Port() int
	ShutdownTimeout() time.Duration
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
