package model

import "errors"

var (
	ErrOrderNotFound             = errors.New("order not found")
	ErrInvalidPaymentMethod      = errors.New("payment method must be specified")
	ErrUnknownPaymentMethod      = errors.New("unknown payment method")
	ErrPaymentClient             = errors.New("payment client error")
	ErrUnexpectedPaymentResponse = errors.New("unexpected payment response type")
	ErrUpdateOrderFailed         = errors.New("failed to update order")
	ErrInvalidTransactionUUID    = errors.New("invalid transaction uuid")
	ErrOrderAlreadyPaid          = errors.New("order already paid")
	ErrOrderCancelled            = errors.New("order is cancelled")
	ErrInvalidOrderStatus        = errors.New("invalid order status for payment")
)

// Validation and request errors for service layer
var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInvalidOrderUUID  = errors.New("invalid order uuid")
	ErrInvalidUserUUID   = errors.New("invalid user uuid")
	ErrEmptyPartUUIDs    = errors.New("part_uuids is required")
	ErrNoFilterProvided  = errors.New("either part_uuids or parts_filter must be provided")
	ErrMutuallyExclusive = errors.New("part_uuids and parts_filter are mutually exclusive")
	ErrPartsNotFound     = errors.New("no parts found by filter")
)
