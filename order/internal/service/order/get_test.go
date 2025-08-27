package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

func (s *ServiceSuite) TestGetOrder_Success() {
	id := uuid.New()
	s.orderRepository.EXPECT().GetOrder(mock.Anything, id).Return(model.Order{OrderUUID: id}, nil)
	res, err := s.service.GetOrder(context.Background(), id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), id, res.OrderUUID)
}

func (s *ServiceSuite) TestGetOrder_NotFound() {
	id := uuid.New()
	s.orderRepository.EXPECT().GetOrder(mock.Anything, id).Return(model.Order{}, model.ErrOrderNotFound)
	_, err := s.service.GetOrder(context.Background(), id)
	require.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}
