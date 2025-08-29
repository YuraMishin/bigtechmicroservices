package order

import (
	"errors"

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
		return nil, errors.New("orderRepository is nil")
	}

	if inventoryClient == nil {
		return nil, errors.New("inventoryClient is nil")
	}

	if paymentClient == nil {
		return nil, errors.New("paymentClient is nil")
	}

	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}, nil
}
