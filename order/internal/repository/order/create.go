package order

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
)

func (r *repository) CreateNewOrder(ctx context.Context, model model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order := converter.ModelToOrder(model)
	r.data[order.OrderUUID.String()] = order
}
