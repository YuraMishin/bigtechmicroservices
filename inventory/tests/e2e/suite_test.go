package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/YuraMishin/bigtechmicroservices/platform/pkg/logger"
	"github.com/YuraMishin/bigtechmicroservices/platform/pkg/testcontainers"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inventory Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Пытаемся загрузить .env; если нет — используем значения по умолчанию
	envPath := filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env")
	if _, statErr := os.Stat(envPath); statErr == nil {
		if envVars, readErr := godotenv.Read(envPath); readErr == nil {
			for key, value := range envVars {
				_ = os.Setenv(key, value)
			}
		} else {
			logger.Warn(suiteCtx, "Не удалось загрузить .env файл, применяю значения по умолчанию", zap.Error(readErr))
		}
	} else {
		logger.Warn(suiteCtx, "Файл .env не найден, применяю значения по умолчанию", zap.Error(statErr))
	}

	// Гарантируем необходимые переменные окружения
	setDefaultEnv("GRPC_PORT", "50051")
	setDefaultEnv("LOGGER_LEVEL", loggerLevelValue)
	setDefaultEnv("LOGGER_AS_JSON", "false")

	setDefaultEnv(testcontainers.MongoImageNameKey, "mongo:8.0")
	setDefaultEnv(testcontainers.MongoDatabaseKey, "inventory")
	setDefaultEnv(testcontainers.MongoUsernameKey, "root")
	setDefaultEnv(testcontainers.MongoPasswordKey, "root")
	setDefaultEnv(testcontainers.MongoAuthDBKey, "admin")
	setDefaultEnv(testcontainers.MongoPortKey, testcontainers.MongoPort)

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})

func setDefaultEnv(key, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
