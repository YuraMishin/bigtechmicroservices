package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	response, err := a.paymentService.PayOrder(ctx, req)
	if err != nil {
		log.Printf("Error in PayOrder: %v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return response, nil
}
