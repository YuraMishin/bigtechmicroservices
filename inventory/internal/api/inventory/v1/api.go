package v1

import (
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/service"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	inventoryService service.PartService
}

func NewAPI(inventoryService service.PartService) *api {
	if inventoryService == nil {
		panic("internal error")
	}
	return &api{
		inventoryService: inventoryService,
	}
}
