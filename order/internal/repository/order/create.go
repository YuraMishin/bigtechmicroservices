package order

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, model model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := converter.ToRepoOrder(model)
	r.data[order.OrderUUID.String()] = order

	return model, nil
}
