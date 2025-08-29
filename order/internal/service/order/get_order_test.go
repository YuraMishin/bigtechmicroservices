package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

func (s *ServiceSuite) TestGetOrder_NotFound_OnNilUUID() {
	s.orderRepository.EXPECT().GetOrder(mock.Anything, uuid.Nil).Return(model.Order{}, model.ErrOrderNotFound)
	_, err := s.service.GetOrder(context.Background(), uuid.Nil)
	require.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestGetOrder_Success() {
	id := uuid.New()
	expected := model.Order{OrderUUID: id}
	s.orderRepository.EXPECT().GetOrder(mock.Anything, id).Return(expected, nil)
	got, err := s.service.GetOrder(context.Background(), id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected, got)
}
