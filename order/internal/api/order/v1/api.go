package v1

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/service"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

type api struct {
	orderService  service.OrderService
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
}

func (a *api) Close() {
	if a.inventoryConn != nil {
		if err := a.inventoryConn.Close(); err != nil {
			log.Printf("Error closing inventory connection: %v", err)
		}
	}
	if a.paymentConn != nil {
		if err := a.paymentConn.Close(); err != nil {
			log.Printf("Error closing payment connection: %v", err)
		}
	}
}

func (a *api) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	log.Printf("Internal error: %v", err)
	return &orderV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(500),
			Message: orderV1.NewOptString("Internal error"),
		},
	}
}
