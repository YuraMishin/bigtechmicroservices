package part

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository"
	def "github.com/YuraMishin/bigtechmicroservices/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	inventoryRepository repository.PartRepository
}

func NewService(inventoryRepository repository.PartRepository) (*service, error) {
	if inventoryRepository == nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &service{
		inventoryRepository: inventoryRepository,
	}, nil
}
