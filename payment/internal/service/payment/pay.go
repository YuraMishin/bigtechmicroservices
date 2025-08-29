package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/model"
)

func (s *service) PayOrder(ctx context.Context, req model.PayOrderRequest) (model.PayOrderResponse, error) {
	switch req.PaymentMethod {
	case model.PaymentMethodCard, model.PaymentMethodSBP, model.PaymentMethodCreditCard, model.PaymentMethodInvestorMoney:

	case model.PaymentMethodUnspecified:
		return model.PayOrderResponse{}, model.ErrInvalidPaymentMethod

	default:
		return model.PayOrderResponse{}, model.ErrUnknownPaymentMethod
	}

	transactionUUID := uuid.New()

	return model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
