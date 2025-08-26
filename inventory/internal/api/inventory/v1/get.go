package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/converter"
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uuid := req.GetUuid()
	part, err := a.inventoryService.GetPart(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrInvalidRequest) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request")
		}
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", uuid)
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &inventoryV1.GetPartResponse{
		Part: converter.ToProtoPart(part),
	}, nil
}
