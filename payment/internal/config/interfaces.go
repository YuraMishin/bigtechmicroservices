package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentGRPCConfig interface {
	Port() int
	ShutdownTimeout() time.Duration
}
