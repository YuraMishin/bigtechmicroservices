package converter

import (
	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func ToOrderDto(order model.Order) *orderV1.OrderDto {
	return &orderV1.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		Status:          order.Status,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
	}
}
