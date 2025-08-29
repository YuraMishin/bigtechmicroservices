package v1

import (
	"context"
	"errors"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	// Конвертируем API-модель в доменную модель
	paymentMethod, err := a.convertPaymentMethod(req.PaymentMethod)
	if err != nil {
		return &orderV1.BadRequestError{Code: 400, Message: err.Error()}, nil
	}

	paymentRequest := model.PaymentRequest{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: paymentMethod,
	}

	paymentResult, err := a.orderService.PayOrder(ctx, paymentRequest)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderUUID), errors.Is(err, model.ErrInvalidPaymentMethod):
			return &orderV1.BadRequestError{Code: 400, Message: err.Error()}, nil

		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "Order not found"}, nil

		case errors.Is(err, model.ErrOrderAlreadyPaid), errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.BadRequestError{Code: 400, Message: err.Error()}, nil

		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal error"}, nil
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: paymentResult.TransactionUUID,
	}, nil
}

func (a *api) convertPaymentMethod(paymentMethod orderV1.PayOrderRequestPaymentMethod) (model.PaymentMethod, error) {
	switch paymentMethod {
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
