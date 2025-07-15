package v1

import (
	"context"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, err := a.orderService.GetOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}
	orderRes, err := a.orderService.PayOrder(ctx, order, req)
	if err != nil {
		return nil, err
	}
	return orderRes, nil
}
