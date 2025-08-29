package v1

import (
	"context"
	"errors"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	res, err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderUUID):
			return &orderV1.BadRequestError{Code: 400, Message: "Invalid order UUID"}, nil

		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "Order not found"}, nil

		case errors.Is(err, model.ErrOrderAlreadyPaid), errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.Conflict{Code: 409, Message: err.Error()}, nil

		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal error"}, nil
		}
	}

	return res, nil
}
