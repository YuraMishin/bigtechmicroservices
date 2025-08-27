package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/model"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (s *service) PayOrder(ctx context.Context, in *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	switch in.PaymentMethod {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		return nil, model.ErrInvalidPaymentMethod

	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:

	default:
		return nil, model.ErrUnknownPaymentMethod
	}

	transactionUUID := uuid.New().String()

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
