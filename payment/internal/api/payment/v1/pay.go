package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/model"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Конвертируем gRPC модель в доменную модель
	paymentRequest := model.PayOrderRequest{
		OrderUUID:     req.OrderUuid,
		UserUUID:      req.UserUuid,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod.String()),
	}

	paymentResponse, err := a.paymentService.PayOrder(ctx, paymentRequest)
	if err != nil {
		log.Printf("Error in PayOrder: %v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: paymentResponse.TransactionUUID.String(),
	}, nil
}
