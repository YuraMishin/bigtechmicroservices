package model

import "errors"

var (
	ErrInvalidPaymentMethod = errors.New("payment method must be specified")
	ErrUnknownPaymentMethod = errors.New("unknown payment method")
)
