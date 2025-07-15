package v1

import (
	"context"

	"github.com/google/uuid"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrderByUUID(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
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
	uuidNoContent, err := a.orderService.CancelOrderByUUID(ctx, order)
	if err != nil {
		return nil, err
	}
	return uuidNoContent, nil
}
