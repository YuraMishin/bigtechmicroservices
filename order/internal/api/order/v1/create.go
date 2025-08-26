package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) CreateNewOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateNewOrderRes, error) {
	// Basic validation of required fields
	if req == nil || req.UserUUID == uuid.Nil || len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "user_uuid and part_uuids are required",
		}, nil
	}

	// Convert to service-layer model
	newOrder := model.NewOrder{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	}
	_ = newOrder // currently service expects OpenAPI request; conversion is prepared

	result, err := a.orderService.CreateNewOrder(ctx, req)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}
	return result, nil
}
