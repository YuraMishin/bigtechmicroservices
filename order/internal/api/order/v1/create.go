package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req == nil {
		return &orderV1.BadRequestError{Code: 400, Message: "request is required"}, nil
	}
	if _, err := uuid.Parse(req.UserUUID.String()); err != nil || req.UserUUID == uuid.Nil {
		return &orderV1.BadRequestError{Code: 400, Message: "invalid user_uuid"}, nil
	}
	// Allow either part_uuids or parts_filter; mutual exclusivity handled in service.
	if len(req.PartUuids) > 0 {
		for _, id := range req.PartUuids {
			if _, err := uuid.Parse(id.String()); err != nil || id == uuid.Nil {
				return &orderV1.BadRequestError{Code: 400, Message: "invalid part_uuid"}, nil
			}
		}
	}

	result, err := a.orderService.CreateOrder(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidRequest),
			errors.Is(err, model.ErrInvalidUserUUID),
			errors.Is(err, model.ErrNoFilterProvided),
			errors.Is(err, model.ErrMutuallyExclusive),
			errors.Is(err, model.ErrPartsNotFound):
			return &orderV1.BadRequestError{Code: 400, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal error"}, nil
		}
	}
	return result, nil
}
