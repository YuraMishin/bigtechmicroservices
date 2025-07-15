package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	repoConverter "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
)

func (r *repository) GetOrderByUUID(_ context.Context, orderUUID uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.data[orderUUID.String()]
	if !exists {
		return model.Order{}, model.ErrOrderNotFound
	}
	return repoConverter.OrderToModel(order), nil
}
