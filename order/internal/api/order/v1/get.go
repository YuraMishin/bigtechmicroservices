package v1

import (
	"context"
	"errors"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/converter"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "Order not found"}, nil

		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal error"}, nil
		}
	}

	return converter.ToOrderDto(order), nil
}
