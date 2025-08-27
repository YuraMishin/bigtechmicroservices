package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/service"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	inventoryService service.PartService
}

func NewAPI(inventoryService service.PartService) (*api, error) {
	if inventoryService == nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &api{
		inventoryService: inventoryService,
	}, nil
}
