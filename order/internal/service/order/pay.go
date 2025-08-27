package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *service) PayOrder(ctx context.Context, orderUUID uuid.UUID, req *orderV1.PayOrderRequest) (orderV1.PayOrderRes, error) {
	if _, err := uuid.Parse(orderUUID.String()); err != nil {
		return nil, model.ErrInvalidOrderUUID
	}
	if req == nil {
		return nil, model.ErrInvalidRequest
	}

	pm, err := s.validateAndMapPaymentMethod(req)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	if order.Status == orderV1.OrderDtoStatusPAID {
		return nil, model.ErrOrderAlreadyPaid
	}
	if order.Status == orderV1.OrderDtoStatusCANCELLED {
		return nil, model.ErrOrderCancelled
	}

	payRes, err := s.paymentClient.PayOrder(ctx, order.OrderUUID.String(), order.UserUUID.String(), pm.String())
	if err != nil {
		return nil, model.ErrPaymentClient
	}

	order.TransactionUUID = payRes.TransactionUUID
	order.PaymentMethod = orderV1.OrderDtoPaymentMethod(pm.String())
	order.Status = orderV1.OrderDtoStatusPAID
	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return nil, model.ErrUpdateOrderFailed
	}

	return &orderV1.PayOrderResponse{TransactionUUID: payRes.TransactionUUID}, nil
}

func (s *service) validateAndMapPaymentMethod(req *orderV1.PayOrderRequest) (model.PaymentMethod, error) {
	switch req.PaymentMethod {
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		return model.PaymentMethodCard, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		return model.PaymentMethodSBP, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		return model.PaymentMethodCreditCard, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethodUnspecified, model.ErrInvalidPaymentMethod
	}
}
