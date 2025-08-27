package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestCancelOrder_Success() {
	id := uuid.New()
	ord := model.Order{OrderUUID: id, Status: orderV1.OrderDtoStatusPENDINGPAYMENT}
	s.orderRepository.EXPECT().GetOrder(mock.Anything, id).Return(ord, nil)
	s.orderRepository.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)
	res, err := s.service.CancelOrder(context.Background(), id)
	require.NoError(s.T(), err)
	_, ok := res.(*orderV1.CancelOrderNoContent)
	require.True(s.T(), ok)
}

func (s *ServiceSuite) TestCancelOrder_AlreadyCancelled() {
	id := uuid.New()
	ord := model.Order{OrderUUID: id, Status: orderV1.OrderDtoStatusCANCELLED}
	s.orderRepository.EXPECT().GetOrder(mock.Anything, id).Return(ord, nil)
	res, err := s.service.CancelOrder(context.Background(), id)
	require.ErrorIs(s.T(), err, model.ErrOrderCancelled)
	require.Nil(s.T(), res)
}
