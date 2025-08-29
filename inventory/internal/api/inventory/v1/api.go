package v1

import (
	"errors"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/service"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	inventoryService service.PartService
}

func NewAPI(inventoryService service.PartService) (*api, error) {
	if inventoryService == nil {
		return nil, errors.New("inventoryService is nil")
	}

	return &api{
		inventoryService: inventoryService,
	}, nil
}
