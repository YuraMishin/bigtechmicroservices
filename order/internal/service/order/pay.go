package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *service) PayOrder(ctx context.Context, paymentRequest model.PaymentRequest) (model.PaymentResult, error) {
	if _, err := uuid.Parse(paymentRequest.OrderUUID.String()); err != nil {
		return model.PaymentResult{}, model.ErrInvalidOrderUUID
	}

	order, err := s.orderRepository.GetOrder(ctx, paymentRequest.OrderUUID)
	if err != nil {
		return model.PaymentResult{}, err
	}

	if order.Status == orderV1.OrderDtoStatusPAID {
		return model.PaymentResult{}, model.ErrOrderAlreadyPaid
	}
	if order.Status == orderV1.OrderDtoStatusCANCELLED {
		return model.PaymentResult{}, model.ErrOrderCancelled
	}

	payRes, err := s.paymentClient.PayOrder(ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentRequest.PaymentMethod.String())
	if err != nil {
		return model.PaymentResult{}, model.ErrPaymentClient
	}

	if _, err := uuid.Parse(payRes.TransactionUUID.String()); err != nil {
		return model.PaymentResult{}, model.ErrInvalidTransactionUUID
	}

	order.TransactionUUID = payRes.TransactionUUID

	order.PaymentMethod = orderV1.OrderDtoPaymentMethod(paymentRequest.PaymentMethod.String())
	order.Status = orderV1.OrderDtoStatusPAID
	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return model.PaymentResult{}, model.ErrUpdateOrderFailed
	}

	return payRes, nil
}
