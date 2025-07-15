package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/converter"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	if params.OrderUUID == uuid.Nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid order UUID",
		}, nil
	}
	order, err := a.orderService.GetOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}
	orderDto := converter.OrderModelToOrderDto(order)
	return orderDto, nil
}
