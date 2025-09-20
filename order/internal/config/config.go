package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger          LoggerConfig
	OrderHTTP       OrderHTTPConfig
	OrderGRPC       OrderGRPCConfig
	OrderPostgresql OrderPostgresConfig
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

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	orderGRPCCfg, err := env.NewOrderGRPCConfig()
	if err != nil {
		return err
	}

	orderPostgresql, err := env.NewOrderPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:          loggerCfg,
		OrderHTTP:       orderHTTPCfg,
		OrderGRPC:       orderGRPCCfg,
		OrderPostgresql: orderPostgresql,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
