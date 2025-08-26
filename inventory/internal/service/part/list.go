package part

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
