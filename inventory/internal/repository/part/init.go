package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
)

func (r *repository) InitializeCollection(ctx context.Context) error {
	if err := r.createIndexes(ctx); err != nil {
		return err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count == 0 {
		if err := r.populateSampleData(ctx); err != nil {
			return err
		}
		log.Println("✅ MongoDB collection initialized with sample data")
	} else {
		log.Printf("📊 MongoDB collection contains %d documents", count)
	}

	return nil
}

func (r *repository) createIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uuid_unique"),
		},
		{
			Keys:    bson.D{{Key: "category", Value: 1}},
			Options: options.Index().SetName("category_idx"),
		},
		{
			Keys:    bson.D{{Key: "manufacturer.country", Value: 1}},
			Options: options.Index().SetName("manufacturer_country_idx"),
		},
		{
			Keys:    bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetName("tags_idx"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetName("name_idx"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *repository) populateSampleData(ctx context.Context) error {
	sampleParts := []repoModel.Part{
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Main Engine Alpha",
			Description:   "High-performance rocket engine for primary propulsion",
			Price:         75000.0,
			StockQuantity: 5,
			Category:      repoModel.CategoryEngine,
			Dimensions: &repoModel.Dimensions{
				Length: 3.5,
				Width:  1.2,
				Height: 2.1,
				Weight: 850.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceTech USA",
				Country: "USA",
				Website: "https://spacetech-usa.com",
			},
			Tags: []string{"propulsion", "thrust", "main-engine"},
			Metadata: map[string]*repoModel.Value{
				"thrust": {
					DoubleValue: func() *float64 { v := 1500000.0; return &v }(),
				},
				"fuel_type": {
					StringValue: func() *string { v := "liquid_hydrogen"; return &v }(),
				},
			},
			CreatedAt: 1640995200, // 2022-01-01
			UpdatedAt: 1640995200,
		},
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Fuel Tank Beta",
			Description:   "Cryogenic fuel storage tank",
			Price:         25000.0,
			StockQuantity: 8,
			Category:      repoModel.CategoryFuel,
			Dimensions: &repoModel.Dimensions{
				Length: 4.0,
				Width:  2.0,
				Height: 2.5,
				Weight: 1200.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Deutsche SpaceTech",
				Country: "Germany",
				Website: "https://deutsche-spacetech.de",
			},
			Tags: []string{"storage", "cryogenic", "fuel"},
			Metadata: map[string]*repoModel.Value{
				"capacity": {
					DoubleValue: func() *float64 { v := 5000.0; return &v }(),
				},
				"temperature": {
					DoubleValue: func() *float64 { v := -253.0; return &v }(),
				},
			},
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		},
	}

	var docs []interface{}
	for _, part := range sampleParts {
		docs = append(docs, part)
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}
