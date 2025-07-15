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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/YuraMishin/bigtechmicroservices/order/internal/api/order/v1"
	inventoryClientV1 "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/order"
	orderService "github.com/YuraMishin/bigtechmicroservices/order/internal/service/order"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	inventoryGrpcPort = 50051
	paymentGrpcPort   = 50052
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	log.Printf("Connecting to InventoryService on localhost:%d", inventoryGrpcPort)
	inventoryConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", inventoryGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}
	log.Printf("Successfully connected to InventoryService")

	log.Printf("Connecting to PaymentService on localhost:%d", paymentGrpcPort)
	paymentConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", paymentGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}
	log.Printf("Successfully connected to PaymentService")

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	inventoryClientAdapter := inventoryClientV1.NewClient(inventoryClient)
	paymentClientAdapter := paymentClientV1.NewClient(paymentClient)

	repo := orderRepository.NewRepository()
	service := orderService.NewService(repo, inventoryClientAdapter, paymentClientAdapter)
	api := orderV1API.NewAPI(inventoryConn, paymentConn, service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		api.Close()
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	defer api.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v", err)
	}

	log.Println("✅  Сервер остановлен")
}
