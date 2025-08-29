package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrder_WithPartUUIDs_Success() {
	p1 := model.Part{UUID: uuid.New().String(), Price: 10}
	p2 := model.Part{UUID: uuid.New().String(), Price: 12.5}

	s.inventoryClient.EXPECT().ListParts(mock.Anything, mock.Anything).Return([]model.Part{p1, p2}, nil)
	s.orderRepository.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(model.Order{}, nil)

	req := model.Order{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{uuid.MustParse(p1.UUID), uuid.MustParse(p2.UUID)},
	}
	res, err := s.service.CreateOrder(context.Background(), req)
	require.NoError(s.T(), err)
	require.NotEqual(s.T(), uuid.Nil, res.OrderUUID)
	require.Greater(s.T(), res.TotalPrice, float32(0))
}

func (s *ServiceSuite) TestCreateOrder_EmptyPartUUIDs_Error() {
	req := model.Order{UserUUID: uuid.New(), PartUUIDs: []uuid.UUID{}}
	res, err := s.service.CreateOrder(context.Background(), req)
	require.ErrorIs(s.T(), err, model.ErrEmptyPartUUIDs)
	require.Equal(s.T(), model.Order{}, res)
}

func (s *ServiceSuite) TestCreateOrder_PartsNotFound_Error() {
	s.inventoryClient.EXPECT().ListParts(mock.Anything, mock.Anything).Return([]model.Part{}, nil)
	req := model.Order{UserUUID: uuid.New(), PartUUIDs: []uuid.UUID{uuid.New()}}
	_, err := s.service.CreateOrder(context.Background(), req)
	require.ErrorIs(s.T(), err, model.ErrPartsNotFound)
}
