package v1

import (
	"context"

	"github.com/google/uuid"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil || params.OrderUUID == uuid.Nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid order UUID",
		}, nil
	}
	res, err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}
	return res, nil
}
