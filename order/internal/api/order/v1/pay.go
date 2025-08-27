package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return &orderV1.BadRequestError{Code: 400, Message: "Invalid order UUID"}, nil
	}
	res, err := a.orderService.PayOrder(ctx, params.OrderUUID, req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderUUID), errors.Is(err, model.ErrInvalidRequest), errors.Is(err, model.ErrInvalidPaymentMethod):
			return &orderV1.BadRequestError{Code: 400, Message: err.Error()}, nil
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "Order not found"}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid), errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.BadRequestError{Code: 400, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal error"}, nil
		}
	}
	return res, nil
}
