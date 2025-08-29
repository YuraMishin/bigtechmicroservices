package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	id := uuid.New()
	user := uuid.New()
	ord := model.Order{OrderUUID: id, UserUUID: user, Status: orderV1.OrderDtoStatusPENDINGPAYMENT}
	s.orderRepository.EXPECT().GetOrder(mock.Anything, id).Return(ord, nil)
	trx := uuid.New()
	s.paymentClient.EXPECT().PayOrder(mock.Anything, id.String(), user.String(), model.PaymentMethodCard.String()).Return(model.PaymentResult{TransactionUUID: trx}, nil)
	s.orderRepository.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

	res, err := s.service.PayOrder(context.Background(), model.PaymentRequest{OrderUUID: id, PaymentMethod: model.PaymentMethodCard})
	require.NoError(s.T(), err)
	require.Equal(s.T(), trx, res.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_AlreadyPaid() {
	id := uuid.New()
	user := uuid.New()
	ord := model.Order{OrderUUID: id, UserUUID: user, Status: orderV1.OrderDtoStatusPAID}
	s.orderRepository.EXPECT().GetOrder(mock.Anything, id).Return(ord, nil)
	res, err := s.service.PayOrder(context.Background(), model.PaymentRequest{OrderUUID: id, PaymentMethod: model.PaymentMethodCard})
	require.ErrorIs(s.T(), err, model.ErrOrderAlreadyPaid)
	require.Equal(s.T(), model.PaymentResult{}, res)
}
