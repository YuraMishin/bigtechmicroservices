package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/YuraMishin/bigtechmicroservices/order/internal/api/order/v1"
	inventoryClientV1 "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc/payment/v1"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/config"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/migrator"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/order"
	orderService "github.com/YuraMishin/bigtechmicroservices/order/internal/service/order"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

// const (
//
//	httpPort          = "8080"
//	inventoryGrpcPort = 50051
//	paymentGrpcPort   = 50052
//	readHeaderTimeout = 5 * time.Second
//	shutdownTimeout   = 10 * time.Second
//
// )
const configPath = "./deploy/compose/order/.env"

func runMigrations(pool *pgxpool.Pool) error {
	migratorRunner := migrator.NewMigrator(
		stdlib.OpenDB(*pool.Config().ConnConfig),
		config.AppConfig().Postgres.MigrationDirectory(),
	)
	if err := migratorRunner.Up(); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func main() {
	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	inventoryConn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", config.AppConfig().Grpc.InventoryGrpcPort()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}

	paymentConn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", config.AppConfig().Grpc.PaymentGrpcPort()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	inventoryClientAdapter := inventoryClientV1.NewClient(inventoryClient)
	paymentClientAdapter := paymentClientV1.NewClient(paymentClient)

	pool, err := pgxpool.New(context.Background(), config.AppConfig().Postgres.URI())
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(pool); err != nil {
		log.Printf("failed to run migrations: %v", err)
		return
	}

	repo, err := order.NewPostgresRepository(pool)
	if err != nil {
		log.Printf("failed to create orderRepository: %v\n", err)
		return
	}

	service, err := orderService.NewService(repo, inventoryClientAdapter, paymentClientAdapter)
	if err != nil {
		log.Printf("failed to create orderService: %v\n", err)
		return
	}

	api, err := orderV1API.NewAPI(inventoryConn, paymentConn, service)
	if err != nil {
		log.Printf("failed to create orderV1API: %v\n", err)
		return
	}

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		api.Close()
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
	}
	defer api.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(config.AppConfig().OrderREST.ShutdownTimeout()))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", config.AppConfig().OrderREST.Port()),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderREST.ReadHeaderTimeout(),
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", config.AppConfig().OrderREST.Port())
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().OrderREST.ShutdownTimeout())
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v", err)
	}
	log.Println("✅  Сервер остановлен")
}
