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

	inventoryV1API "github.com/YuraMishin/bigtechmicroservices/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/part"
	inventoryService "github.com/YuraMishin/bigtechmicroservices/inventory/internal/service/part"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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

	repo := inventoryRepository.NewRepository()

	service, err := inventoryService.NewService(repo)
	if err != nil {
		log.Printf("failed to create inventoryService: %v\n", err)
		return
	}

	api, err := inventoryV1API.NewAPI(service)
	if err != nil {
		log.Printf("failed to create inventoryV1API: %v\n", err)
		return
	}

	inventoryV1.RegisterInventoryServiceServer(s, api)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
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
