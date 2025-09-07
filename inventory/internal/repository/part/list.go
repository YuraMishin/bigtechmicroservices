package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	serviceModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoConverter "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/converter"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter serviceModel.PartsFilter) ([]serviceModel.Part, error) {
	repoFilter := repoConverter.ToRepoPartsFilter(filter)

	mongoFilter := buildMongoFilter(repoFilter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			// Log error but don't return it as the main operation succeeded
			_ = cerr // explicitly ignore the error
		}
	}()

	var repoParts []repoModel.Part
	if err = cursor.All(ctx, &repoParts); err != nil {
		return nil, err
	}

	var serviceParts []serviceModel.Part
	for _, part := range repoParts {
		serviceParts = append(serviceParts, repoConverter.ToModelPart(part))
	}

	return serviceParts, nil
}

func buildMongoFilter(filter repoModel.PartsFilter) bson.M {
	mongoFilter := bson.M{}

	if len(filter.UUIDs) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		mongoFilter["category"] = bson.M{"$in": filter.Categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	return mongoFilter
}
