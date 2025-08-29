package converter

import (
	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	repoModel "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/model"
)

func ToModelOrder(order repoModel.Order) model.Order {
	return model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}

func ToRepoOrder(order model.Order) repoModel.Order {
	return repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}
