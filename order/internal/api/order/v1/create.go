package v1

import (
	"context"
	"errors"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order := model.Order{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	}

	createdOrder, err := a.orderService.CreateOrder(ctx, order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidUserUUID),
			errors.Is(err, model.ErrEmptyPartUUIDs),
			errors.Is(err, model.ErrPartsNotFound):
			return &orderV1.BadRequestError{Code: 400, Message: err.Error()}, nil

		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal error"}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  createdOrder.OrderUUID,
		TotalPrice: createdOrder.TotalPrice,
	}, nil
}
