package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

func (s *service) GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
