package v1

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/converter"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, in *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.PartsFilterToModel(in.GetFilter()))
	if err != nil {
		return nil, err
	}
	protoParts := make([]*inventoryV1.Part, len(parts))
	for i, p := range parts {
		protoParts[i] = converter.PartToProto(p)
	}
	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
