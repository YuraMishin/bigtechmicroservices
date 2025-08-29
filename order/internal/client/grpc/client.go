package grpc

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (model.PaymentResult, error)
}
