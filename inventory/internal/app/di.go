package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryV1API "github.com/YuraMishin/bigtechmicroservices/inventory/internal/api/inventory/v1"
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/config"
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository"
	inventoryRepository "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/part"
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/service"
	inventoryService "github.com/YuraMishin/bigtechmicroservices/inventory/internal/service/part"
	"github.com/YuraMishin/bigtechmicroservices/platform/pkg/closer"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	inventoryService service.PartService

	inventoryRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	var err error
	if d.inventoryV1API == nil {
		d.inventoryV1API, err = inventoryV1API.NewAPI(d.PartService(ctx))
		if err != nil {
			panic(fmt.Sprintf("failed to create inventoryV1API: %v\n", err))
		}
	}
	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.inventoryService == nil {
		var err error
		d.inventoryService, err = inventoryService.NewService(d.PartRepository(ctx))
		if err != nil {
			panic(fmt.Sprintf("failed to create inventoryService: %v\n", err))
		}
	}

	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(d.MongoDBClient(ctx), config.AppConfig().Mongo.DatabaseName())
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
