package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	if _, err := uuid.Parse(orderUUID.String()); err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
