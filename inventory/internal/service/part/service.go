package part

import (
	"errors"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository"
	def "github.com/YuraMishin/bigtechmicroservices/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	inventoryRepository repository.PartRepository
}

func NewService(inventoryRepository repository.PartRepository) (*service, error) {
	if inventoryRepository == nil {
		return nil, errors.New("inventoryRepository is nil")
	}

	return &service{
		inventoryRepository: inventoryRepository,
	}, nil
}
