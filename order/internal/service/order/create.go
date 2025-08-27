package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *service) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req == nil {
		return nil, model.ErrInvalidRequest
	}
	if _, err := uuid.Parse(req.UserUUID.String()); err != nil || req.UserUUID == uuid.Nil {
		return nil, model.ErrInvalidUserUUID
	}

	var filter model.PartsFilter
	if len(req.PartUuids) > 0 && req.PartsFilter.IsSet() {
		return nil, model.ErrMutuallyExclusive
	}

	hasPartUUIDs := len(req.PartUuids) > 0
	hasPartsFilter := req.PartsFilter.IsSet()

	switch {
	case hasPartsFilter:
		pf := req.PartsFilter.Value
		uuidStrings := make([]string, 0, len(pf.Uuids))
		for _, id := range pf.Uuids {
			uuidStrings = append(uuidStrings, id.String())
		}
		categories := make([]model.Category, 0, len(pf.Categories))
		for _, c := range pf.Categories {
			categories = append(categories, model.Category(c))
		}
		filter = model.PartsFilter{
			UUIDs:                 uuidStrings,
			Names:                 append([]string{}, pf.Names...),
			Categories:            categories,
			ManufacturerCountries: append([]string{}, pf.ManufacturerCountries...),
			Tags:                  append([]string{}, pf.Tags...),
		}
	case hasPartUUIDs:
		uuidStrings := make([]string, 0, len(req.PartUuids))
		for _, id := range req.PartUuids {
			uuidStrings = append(uuidStrings, id.String())
		}
		filter = model.PartsFilter{UUIDs: uuidStrings}
	default:
		return nil, model.ErrNoFilterProvided
	}

	parts, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, model.ErrPartsNotFound
	}

	var total float64
	resolvedPartUUIDs := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		total += p.Price
		if id, err := uuid.Parse(p.UUID); err == nil {
			resolvedPartUUIDs = append(resolvedPartUUIDs, id)
		}
	}

	newOrder := model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   req.UserUUID,
		PartUUIDs:  resolvedPartUUIDs,
		TotalPrice: float32(total),
		Status:     orderV1.OrderDtoStatusPENDINGPAYMENT,
	}

	_, err = s.orderRepository.CreateOrder(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  newOrder.OrderUUID,
		TotalPrice: newOrder.TotalPrice,
	}, nil
}
