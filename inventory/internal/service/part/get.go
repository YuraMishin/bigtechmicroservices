package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuidStr string) (model.Part, error) {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return model.Part{}, model.ErrInvalidRequest
	}
	part, err := s.inventoryRepository.GetPart(ctx, uuidStr)
	if err != nil {
		return model.Part{}, err
	}
	return part, nil
}
