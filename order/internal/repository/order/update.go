package order

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[order.OrderUUID.String()]; !exists {
		return model.ErrOrderNotFound
	}

	r.data[order.OrderUUID.String()] = converter.ToRepoOrder(order)

	return nil
}
