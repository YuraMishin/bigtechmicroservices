package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) (orderV1.CancelOrderRes, error)
	CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error)
	PayOrder(ctx context.Context, orderUUID uuid.UUID, req *orderV1.PayOrderRequest) (orderV1.PayOrderRes, error)
}
