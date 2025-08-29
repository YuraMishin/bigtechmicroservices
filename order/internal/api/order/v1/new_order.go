package v1

import (
	"errors"

	"google.golang.org/grpc"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/service"
)

func NewAPI(inventoryConn, paymentConn *grpc.ClientConn, orderService service.OrderService) (*api, error) {
	if inventoryConn == nil {
		return nil, errors.New("inventoryConn is nil")
	}

	if paymentConn == nil {
		return nil, errors.New("paymentConn is nil")
	}

	if orderService == nil {
		return nil, errors.New("orderService is nil")
	}

	handler := &api{
		orderService: orderService,
	}

	handler.inventoryConn = inventoryConn
	handler.paymentConn = paymentConn

	return handler, nil
}
