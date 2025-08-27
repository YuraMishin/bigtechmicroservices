package v1

import (
	"context"
	"log"

	clientConverter "github.com/YuraMishin/bigtechmicroservices/order/internal/client/converter"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	generatedInventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	protoFilter := clientConverter.FilterToProto(filter)

	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: protoFilter,
	})
	if err != nil {
		log.Printf("Inventory client: gRPC error: %v", err)
		return nil, err
	}

	return clientConverter.ListToModel(parts.Parts), nil
}
