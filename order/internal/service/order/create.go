package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *service) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	if order.UserUUID == uuid.Nil {
		return model.Order{}, model.ErrInvalidUserUUID
	}

	if len(order.PartUUIDs) == 0 {
		return model.Order{}, model.ErrEmptyPartUUIDs
	}

	uuidStrings := make([]string, 0, len(order.PartUUIDs))
	for _, id := range order.PartUUIDs {
		uuidStrings = append(uuidStrings, id.String())
	}
	filter := model.PartsFilter{UUIDs: uuidStrings}

	parts, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return model.Order{}, err
	}
	if len(parts) == 0 {
		return model.Order{}, model.ErrPartsNotFound
	}

	var total float64
	resolvedPartUUIDs := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		total += p.Price
		if id, err := uuid.Parse(p.UUID); err == nil {
			resolvedPartUUIDs = append(resolvedPartUUIDs, id)
		}
	}

	createdOrder := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        order.UserUUID,
		PartUUIDs:       resolvedPartUUIDs,
		TotalPrice:      float32(total),
		TransactionUUID: uuid.Nil, // Будет заполнено при оплате
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
	}

	_, err = s.orderRepository.CreateOrder(ctx, createdOrder)
	if err != nil {
		return model.Order{}, err
	}

	return createdOrder, nil
}
