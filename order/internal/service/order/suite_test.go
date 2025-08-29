package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc/mocks"
	repoMocks "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	orderRepository *repoMocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	var err error
	s.service, err = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
	if err != nil {
		s.T().Fatalf("failed to create service: %v", err)
	}
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
