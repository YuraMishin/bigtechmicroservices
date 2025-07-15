package order

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, model model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[model.OrderUUID.String()] = converter.ModelToOrder(model)
}
