package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestCreateOrder_WithPartUUIDs_Success() {
	p1 := model.Part{UUID: uuid.New().String(), Price: 10}
	p2 := model.Part{UUID: uuid.New().String(), Price: 12.5}

	s.inventoryClient.EXPECT().ListParts(mock.Anything, mock.Anything).Return([]model.Part{p1, p2}, nil)
	s.orderRepository.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(model.Order{}, nil)

	req := &orderV1.CreateOrderRequest{
		UserUUID:  uuid.New(),
		PartUuids: []uuid.UUID{uuid.MustParse(p1.UUID), uuid.MustParse(p2.UUID)},
	}
	res, err := s.service.CreateOrder(context.Background(), req)
	require.NoError(s.T(), err)
	_, ok := res.(*orderV1.CreateOrderResponse)
	require.True(s.T(), ok)
}

func (s *ServiceSuite) TestCreateOrder_WithFilter_Success() {
	p := model.Part{UUID: uuid.New().String(), Price: 5}
	s.inventoryClient.EXPECT().ListParts(mock.Anything, mock.Anything).Return([]model.Part{p}, nil)
	s.orderRepository.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(model.Order{}, nil)

	pf := orderV1.CreateOrderRequestPartsFilter{Uuids: []uuid.UUID{uuid.MustParse(p.UUID)}}
	req := &orderV1.CreateOrderRequest{UserUUID: uuid.New(), PartsFilter: orderV1.NewOptCreateOrderRequestPartsFilter(pf)}
	res, err := s.service.CreateOrder(context.Background(), req)
	require.NoError(s.T(), err)
	_, ok := res.(*orderV1.CreateOrderResponse)
	require.True(s.T(), ok)
}

func (s *ServiceSuite) TestCreateOrder_Exclusive_Error() {
	p := uuid.New()
	pf := orderV1.CreateOrderRequestPartsFilter{Uuids: []uuid.UUID{p}}
	req := &orderV1.CreateOrderRequest{UserUUID: uuid.New(), PartUuids: []uuid.UUID{p}, PartsFilter: orderV1.NewOptCreateOrderRequestPartsFilter(pf)}
	res, err := s.service.CreateOrder(context.Background(), req)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}
