package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/converter"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, r *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.ToModelPartsFilter(r.GetFilter()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.ToProtoPartList(parts),
	}, nil
}
