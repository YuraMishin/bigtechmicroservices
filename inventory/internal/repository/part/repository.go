package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client, databaseName string) *repository {
	database := client.Database(databaseName)
	collection := database.Collection("parts")

	return &repository{
		client:     client,
		database:   database,
		collection: collection,
	}
}
