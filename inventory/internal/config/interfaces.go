package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Port() int
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
