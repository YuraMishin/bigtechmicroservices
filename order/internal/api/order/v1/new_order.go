package v1

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/service"
)

func NewAPI(inventoryConn, paymentConn *grpc.ClientConn, orderService service.OrderService) (*api, error) {
	if inventoryConn == nil {
		return nil, status.Errorf(codes.InvalidArgument, "inventoryConn is nil")
	}

	if paymentConn == nil {
		return nil, status.Errorf(codes.InvalidArgument, "paymentConn is nil")
	}

	if orderService == nil {
		return nil, status.Errorf(codes.InvalidArgument, "orderService is nil")
	}

	handler := &api{
		orderService: orderService,
	}

	handler.inventoryConn = inventoryConn
	handler.paymentConn = paymentConn

	return handler, nil
}
