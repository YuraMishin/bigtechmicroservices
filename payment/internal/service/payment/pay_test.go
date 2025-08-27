package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/model"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func TestPayOrder_InvalidPaymentMethod(t *testing.T) {
	svc := NewService()
	_, err := svc.PayOrder(context.Background(), &paymentV1.PayOrderRequest{
		OrderUuid:     uuid.NewString(),
		UserUuid:      uuid.NewString(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrInvalidPaymentMethod)
}

func TestPayOrder_UnknownPaymentMethod(t *testing.T) {
	svc := NewService()
	_, err := svc.PayOrder(context.Background(), &paymentV1.PayOrderRequest{
		OrderUuid:     uuid.NewString(),
		UserUuid:      uuid.NewString(),
		PaymentMethod: paymentV1.PaymentMethod(999),
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrUnknownPaymentMethod)
}

func TestPayOrder_Success_Methods(t *testing.T) {
	tests := []struct {
		name   string
		method paymentV1.PaymentMethod
	}{
		{"CARD", paymentV1.PaymentMethod_PAYMENT_METHOD_CARD},
		{"SBP", paymentV1.PaymentMethod_PAYMENT_METHOD_SBP},
		{"CREDIT_CARD", paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD},
		{"INVESTOR_MONEY", paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY},
	}

	svc := NewService()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := svc.PayOrder(context.Background(), &paymentV1.PayOrderRequest{
				OrderUuid:     uuid.NewString(),
				UserUuid:      uuid.NewString(),
				PaymentMethod: tc.method,
			})
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			_, parseErr := uuid.Parse(resp.GetTransactionUuid())
			assert.NoError(t, parseErr)
			assert.NotEmpty(t, resp.GetTransactionUuid())
		})
	}
}
