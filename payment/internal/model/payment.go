package model

import "github.com/google/uuid"

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "PAYMENT_METHOD_UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodSBP           PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodCreditCard    PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

type PayOrderRequest struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID uuid.UUID
}
