package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/model"
)

func TestPayOrder_InvalidPaymentMethod(t *testing.T) {
	svc := NewService()
	_, err := svc.PayOrder(context.Background(), model.PayOrderRequest{
		OrderUUID:     uuid.NewString(),
		UserUUID:      uuid.NewString(),
		PaymentMethod: model.PaymentMethodUnspecified,
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrInvalidPaymentMethod)
}

func TestPayOrder_UnknownPaymentMethod(t *testing.T) {
	svc := NewService()
	_, err := svc.PayOrder(context.Background(), model.PayOrderRequest{
		OrderUUID:     uuid.NewString(),
		UserUUID:      uuid.NewString(),
		PaymentMethod: model.PaymentMethod("UNKNOWN"),
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrUnknownPaymentMethod)
}

func TestPayOrder_Success_Methods(t *testing.T) {
	tests := []struct {
		name   string
		method model.PaymentMethod
	}{
		{"CARD", model.PaymentMethodCard},
		{"SBP", model.PaymentMethodSBP},
		{"CREDIT_CARD", model.PaymentMethodCreditCard},
		{"INVESTOR_MONEY", model.PaymentMethodInvestorMoney},
	}

	svc := NewService()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := svc.PayOrder(context.Background(), model.PayOrderRequest{
				OrderUUID:     uuid.NewString(),
				UserUUID:      uuid.NewString(),
				PaymentMethod: tc.method,
			})
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, resp.TransactionUUID)
			assert.NotEmpty(t, resp.TransactionUUID.String())
		})
	}
}
