package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *service) CreateNewOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateNewOrderRes, error) {
	orderUUID := uuid.New()
	partUUIDs := make([]string, 0, len(req.PartUuids))
	for _, partUUID := range req.PartUuids {
		partUUIDs = append(partUUIDs, partUUID.String())
	}
	filter := model.PartsFilter{
		UUIDs: partUUIDs,
	}
	inventoryResponse, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}
	if len(inventoryResponse) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Some parts not found in inventory",
		}, nil
	}
	var totalPrice float32
	for _, part := range inventoryResponse {
		totalPrice += float32(part.Price)
	}
	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUuids,
		TotalPrice:      totalPrice,
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}
	s.orderRepository.CreateNewOrder(ctx, order)
	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
