package v1

import (
	"context"
	"log"

	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	log.Printf("PayOrder called with request: %+v", req)

	response, err := a.paymentService.PayOrder(ctx, req)
	if err != nil {
		log.Printf("Error in PayOrder: %v", err)
		return nil, err
	}

	log.Printf("PayOrder response: %+v", response)
	return response, nil
}
