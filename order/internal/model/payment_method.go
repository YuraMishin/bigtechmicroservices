package model

type PaymentMethod int

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

// String returns the canonical string representation expected by the payment client.
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
