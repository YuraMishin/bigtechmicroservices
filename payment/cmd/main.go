package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1API "github.com/YuraMishin/bigtechmicroservices/payment/internal/api/payment/v1"
	"github.com/YuraMishin/bigtechmicroservices/payment/internal/config"
	paymentService "github.com/YuraMishin/bigtechmicroservices/payment/internal/service/payment"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConfig().PaymentGRPC.Port()))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	reflection.Register(s)

	service := paymentService.NewService()
	api, err := paymentV1API.NewAPI(service)
	if err != nil {
		log.Printf("failed to create paymentV1API: %v\n", err)
		return
	}

	paymentV1.RegisterPaymentServiceServer(s, api)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", config.AppConfig().PaymentGRPC.Port())
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
