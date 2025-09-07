package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"

	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет тестовую деталь в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now().Unix()

	doc := bson.M{
		"uuid":           partUUID,
		"name":           gofakeit.AppName(),
		"description":    gofakeit.Sentence(10),
		"price":          gofakeit.Price(10, 10000),
		"stock_quantity": int64(gofakeit.Number(1, 100)),
		"category":       int32(1), // CategoryEngine
		"tags":           []string{},
		"metadata":       bson.M{},
		"created_at":     now,
		"updated_at":     now,
	}

	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовую деталь с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, p *inventoryV1.Part) (string, error) {
	partUUID := p.GetUuid()
	if partUUID == "" {
		partUUID = gofakeit.UUID()
	}
	now := time.Now().Unix()

	doc := bson.M{
		"uuid":           partUUID,
		"name":           p.GetName(),
		"description":    p.GetDescription(),
		"price":          p.GetPrice(),
		"stock_quantity": int64(p.GetStockQuantity()),
		"category":       int32(p.GetCategory()),
		"tags":           []string{},
		"metadata":       bson.M{},
		"created_at":     now,
		"updated_at":     now,
	}

	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// ClearInventoryCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearInventoryCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).DeleteMany(ctx, bson.M{})
	return err
}
