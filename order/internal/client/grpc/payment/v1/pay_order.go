package v1

import (
	"context"
	"log"

	"github.com/google/uuid"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
	generatedPaymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (orderV1.PayOrderRes, error) {
	log.Printf("Payment client: calling with orderUUID=%s, userUUID=%s, paymentMethod=%s", orderUUID, userUUID, paymentMethod)

	// Преобразуем строковый paymentMethod в enum
	var paymentMethodEnum generatedPaymentV1.PaymentMethod
	switch paymentMethod {
	case "PAYMENT_METHOD_CARD":
		paymentMethodEnum = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case "PAYMENT_METHOD_SBP":
		paymentMethodEnum = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case "PAYMENT_METHOD_CREDIT_CARD":
		paymentMethodEnum = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "PAYMENT_METHOD_INVESTOR_MONEY":
		paymentMethodEnum = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		paymentMethodEnum = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	request := &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentMethodEnum,
	}

	log.Printf("Payment client: calling gRPC with request %+v", request)
	response, err := c.generatedClient.PayOrder(ctx, request)
	if err != nil {
		log.Printf("Payment client: gRPC error: %v", err)
		return nil, err
	}
	log.Printf("Payment client: gRPC returned transaction UUID: %s", response.TransactionUuid)

	// Конвертация в тип OpenAPI-модели (uuid.UUID). Ошибки парсинга не прокидываем из клиента.
	var txUUID uuid.UUID
	if parsed, err := uuid.Parse(response.TransactionUuid); err == nil {
		txUUID = parsed
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: txUUID,
	}, nil
}
