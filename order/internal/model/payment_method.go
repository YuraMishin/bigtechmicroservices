package model

import "github.com/google/uuid"

type PaymentMethod int

const (
	PaymentMethodUnspecified   PaymentMethod = 0
	PaymentMethodCard          PaymentMethod = 1
	PaymentMethodSBP           PaymentMethod = 2
	PaymentMethodCreditCard    PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)

func (p PaymentMethod) String() string {
	switch p {
	case PaymentMethodCard:
		return "PAYMENT_METHOD_CARD"
	case PaymentMethodSBP:
		return "PAYMENT_METHOD_SBP"
	case PaymentMethodCreditCard:
		return "PAYMENT_METHOD_CREDIT_CARD"
	case PaymentMethodInvestorMoney:
		return "PAYMENT_METHOD_INVESTOR_MONEY"
	default:
		return "PAYMENT_METHOD_UNSPECIFIED"
	}
}

type PaymentRequest struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethod
}
