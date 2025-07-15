package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

type OrderService interface {
	GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	CancelOrderByUUID(ctx context.Context, order model.Order) (orderV1.CancelOrderByUUIDRes, error)
	CreateNewOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateNewOrderRes, error)
	PayOrder(ctx context.Context, order model.Order, req *orderV1.PayOrderRequest) (orderV1.PayOrderRes, error)
}
