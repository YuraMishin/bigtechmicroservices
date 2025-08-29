package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (model.PaymentResult, error) {
	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[paymentMethod]),
	}

	resp, err := c.generatedClient.PayOrder(ctx, req)
	if err != nil {
		return model.PaymentResult{}, err
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return model.PaymentResult{}, err
	}

	return model.PaymentResult{
		TransactionUUID: transactionUUID,
	}, nil
}
