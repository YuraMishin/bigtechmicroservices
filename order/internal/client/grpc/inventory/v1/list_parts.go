package v1

import (
	"context"
	"log"

	clientConverter "github.com/YuraMishin/bigtechmicroservices/order/internal/client/converter"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	generatedInventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	log.Printf("Inventory client: calling gRPC with filter %+v", filter)
	protoFilter := clientConverter.PartsFilterToProto(filter)
	log.Printf("Inventory client: converted to proto filter %+v", protoFilter)

	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: protoFilter,
	})
	if err != nil {
		log.Printf("Inventory client: gRPC error: %v", err)
		return nil, err
	}
	log.Printf("Inventory client: gRPC returned %d parts", len(parts.Parts))

	result := clientConverter.PartListToModel(parts.Parts)
	log.Printf("Inventory client: converted to %d model parts", len(result))
	return result, nil
}
