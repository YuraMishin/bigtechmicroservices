package order

import (
	"context"
	"log"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *service) CancelOrderByUUID(ctx context.Context, order model.Order) (orderV1.CancelOrderByUUIDRes, error) {
	switch order.Status {
	case orderV1.OrderDtoStatusPENDINGPAYMENT:
		order.Status = orderV1.OrderDtoStatusCANCELLED
		s.orderRepository.UpdateOrder(ctx, order)
		log.Printf("Order %s cancelled successfully", order.OrderUUID.String())
		return &orderV1.CancelOrderByUUIDNoContent{}, nil
	case orderV1.OrderDtoStatusPAID:
		return &orderV1.Conflict{
			Code:    409,
			Message: "Order is already paid and cannot be cancelled",
		}, nil
	case orderV1.OrderDtoStatusCANCELLED:
		return &orderV1.CancelOrderByUUIDNoContent{}, nil
	default:
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid order status",
		}, nil
	}
}
