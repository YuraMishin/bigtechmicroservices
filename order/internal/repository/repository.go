package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	UpdateOrder(ctx context.Context, order model.Order) error
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
}
