package model

import (
	"github.com/google/uuid"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float32
	TransactionUUID uuid.UUID
	PaymentMethod   orderV1.OrderDtoPaymentMethod
	Status          orderV1.OrderDtoStatus
}
