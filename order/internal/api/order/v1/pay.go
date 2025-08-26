package v1

import (
	"context"

	"github.com/google/uuid"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil || params.OrderUUID == uuid.Nil {
		return &orderV1.BadRequestError{Code: 400, Message: "Invalid order UUID"}, nil
	}
	res, err := a.orderService.PayOrderByUUID(ctx, params.OrderUUID, req)
	if err != nil {
		return &orderV1.InternalServerError{Code: 500, Message: "Internal error"}, nil
	}
	return res, nil
}
