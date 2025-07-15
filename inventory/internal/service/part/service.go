package part

import (
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository"
	def "github.com/YuraMishin/bigtechmicroservices/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	inventoryRepository repository.PartRepository
}

func NewService(inventoryRepository repository.PartRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
