package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/YuraMishin/bigtechmicroservices/inventory/internal/api/inventory/v1"
	mongoRepository "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/part"
	inventoryService "github.com/YuraMishin/bigtechmicroservices/inventory/internal/service/part"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
)

func getMongoClient() (*mongo.Client, error) {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDatabase := os.Getenv("MONGO_DATABASE")

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", mongoUser, mongoPassword, mongoHost, mongoPort, mongoDatabase)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load .env, proceeding with environment variables")
	}

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

	client, err := getMongoClient()
	if err != nil {
		log.Printf("failed to connect to mongo: %v\n", err)
		return
	}

	defer func() {
		if cerr := client.Disconnect(context.Background()); cerr != nil {
			log.Printf("failed to disconnect from mongo: %v\n", cerr)
		}
	}()

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Printf("failed to ping mongo: %v\n", err)
		return
	}

	repo := mongoRepository.NewRepository(client, os.Getenv("MONGO_DATABASE"))

	if err := repo.InitializeCollection(context.Background()); err != nil {
		log.Printf("failed to initialize collection: %v\n", err)
		return
	}

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
