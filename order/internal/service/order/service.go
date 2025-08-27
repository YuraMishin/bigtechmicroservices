package order

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository"
	def "github.com/YuraMishin/bigtechmicroservices/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) (*service, error) {
	if orderRepository == nil {
		return nil, status.Errorf(codes.InvalidArgument, "orderRepository is nil")
	}

	if inventoryClient == nil {
		return nil, status.Errorf(codes.InvalidArgument, "inventoryClient is nil")
	}

	if paymentClient == nil {
		return nil, status.Errorf(codes.InvalidArgument, "paymentClient is nil")
	}

	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}, nil
}
