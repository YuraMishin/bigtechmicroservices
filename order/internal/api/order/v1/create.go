package v1

import (
	"context"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) CreateNewOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateNewOrderRes, error) {
	result, err := a.orderService.CreateNewOrder(ctx, req)
	return result, err
}
