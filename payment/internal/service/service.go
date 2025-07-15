package service

import (
	"context"

	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

type PaymentService interface {
	PayOrder(ctx context.Context, in *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error)
}
