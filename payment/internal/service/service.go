package service

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, req model.PayOrderRequest) (model.PayOrderResponse, error)
}
