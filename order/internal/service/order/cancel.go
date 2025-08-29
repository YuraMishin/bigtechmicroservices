package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) (orderV1.CancelOrderRes, error) {
	if _, err := uuid.Parse(orderUUID.String()); err != nil || orderUUID == uuid.Nil {
		return nil, model.ErrInvalidOrderUUID
	}

	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	if order.Status == orderV1.OrderDtoStatusCANCELLED {
		return nil, model.ErrOrderCancelled
	}
	if order.Status == orderV1.OrderDtoStatusPAID {
		return nil, model.ErrOrderAlreadyPaid
	}

	order.Status = orderV1.OrderDtoStatusCANCELLED
	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
