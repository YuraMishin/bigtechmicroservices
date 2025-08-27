package v1

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (model.PaymentResult, error) {
	return model.PaymentResult{}, nil
}
