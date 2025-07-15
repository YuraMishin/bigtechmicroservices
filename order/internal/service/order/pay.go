package order

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (s *service) PayOrder(ctx context.Context, order model.Order, req *orderV1.PayOrderRequest) (orderV1.PayOrderRes, error) {
	var paymentMethod paymentV1.PaymentMethod
	switch req.PaymentMethod {
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED:
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Payment method must be specified",
		}, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Unknown payment method",
		}, nil
	}
	paymentRequest := &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	}
	var paymentMethodString string
	switch paymentMethod {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		paymentMethodString = "PAYMENT_METHOD_CARD"
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		paymentMethodString = "PAYMENT_METHOD_SBP"
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		paymentMethodString = "PAYMENT_METHOD_CREDIT_CARD"
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		paymentMethodString = "PAYMENT_METHOD_INVESTOR_MONEY"
	default:
		paymentMethodString = "PAYMENT_METHOD_UNSPECIFIED"
	}
	paymentResponse, err := s.paymentClient.PayOrder(ctx, paymentRequest.OrderUuid, paymentRequest.UserUuid, paymentMethodString)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}
	paymentResult, ok := paymentResponse.(*orderV1.PayOrderResponse)
	if !ok {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}
	transactionUUID := paymentResult.TransactionUUID
	var orderPaymentMethod orderV1.OrderDtoPaymentMethod
	switch req.PaymentMethod {
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY
	default:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
	}
	order.Status = orderV1.OrderDtoStatusPAID
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = orderPaymentMethod
	s.orderRepository.UpdateOrder(ctx, order)
	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
