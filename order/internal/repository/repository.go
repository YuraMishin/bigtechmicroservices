package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

type OrderRepository interface {
	GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	UpdateOrder(ctx context.Context, order model.Order)
	CreateNewOrder(ctx context.Context, order model.Order)
}
