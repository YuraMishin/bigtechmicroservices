package order

import (
	"fmt"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

// parsePaymentMethod конвертирует строку в OrderDtoPaymentMethod
func parsePaymentMethod(s string) (orderV1.OrderDtoPaymentMethod, error) {
	switch s {
	case "PAYMENT_METHOD_UNSPECIFIED":
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED, nil
	case "PAYMENT_METHOD_CARD":
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD, nil
	case "PAYMENT_METHOD_SBP":
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP, nil
	case "PAYMENT_METHOD_CREDIT_CARD":
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD, nil
	case "PAYMENT_METHOD_INVESTOR_MONEY":
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY, nil
	default:
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED, fmt.Errorf("unknown payment method: %s", s)
	}
}

// parseOrderStatus конвертирует строку в OrderDtoStatus
func parseOrderStatus(s string) (orderV1.OrderDtoStatus, error) {
	switch s {
	case "PENDING_PAYMENT":
		return orderV1.OrderDtoStatusPENDINGPAYMENT, nil
	case "PAID":
		return orderV1.OrderDtoStatusPAID, nil
	case "CANCELLED":
		return orderV1.OrderDtoStatusCANCELLED, nil
	default:
		return orderV1.OrderDtoStatusPENDINGPAYMENT, fmt.Errorf("unknown order status: %s", s)
	}
}
