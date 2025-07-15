package v1

import (
	"google.golang.org/grpc"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/service"
)

func NewAPI(inventoryConn, paymentConn *grpc.ClientConn, orderService service.OrderService) *api {
	handler := &api{
		orderService: orderService,
	}
	handler.inventoryConn = inventoryConn
	handler.paymentConn = paymentConn
	return handler
}
