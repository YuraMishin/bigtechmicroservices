package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	partRepository *mocks.PartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.partRepository = mocks.NewPartRepository(s.T())

	s.service, _ = NewService(
		s.partRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
