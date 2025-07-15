package part

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/converter"
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	repoFilter := converter.PartsFilterToRepo(filter)
	parts, err := s.inventoryRepository.ListParts(ctx, repoFilter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
