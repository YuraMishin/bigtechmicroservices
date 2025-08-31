package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger    LoggerConfig
	OrderREST OrderRESTConfig
	Postgres  PostgresConfig
	Grpc      GRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderRESTCfg, err := env.NewOrderRESTConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewGrpcConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:    loggerCfg,
		OrderREST: orderRESTCfg,
		Postgres:  postgresCfg,
		Grpc:      grpcCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
