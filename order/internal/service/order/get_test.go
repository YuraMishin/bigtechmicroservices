package order

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestGetOrderByUUID_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
	}

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(expectedOrder, nil)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOrder, result)
	assert.Equal(s.T(), orderUUID, result.OrderUUID)
	assert.Equal(s.T(), userUUID, result.UserUUID)
	assert.Equal(s.T(), partUUIDs, result.PartUUIDs)
	assert.Equal(s.T(), orderV1.OrderDtoStatusPENDINGPAYMENT, result.Status)
}

func (s *ServiceSuite) TestGetOrderByUUID_OrderNotFound() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(model.Order{}, model.ErrOrderNotFound)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *ServiceSuite) TestGetOrderByUUID_RepositoryError() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	expectedError := errors.New("database connection failed")

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(model.Order{}, expectedError)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *ServiceSuite) TestGetOrderByUUID_WithPaidOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(500, 2000),
		Status:          orderV1.OrderDtoStatusPAID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD,
	}

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(expectedOrder, nil)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOrder, result)
	assert.Equal(s.T(), orderV1.OrderDtoStatusPAID, result.Status)
	assert.Equal(s.T(), transactionUUID, result.TransactionUUID)
	assert.Equal(s.T(), orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD, result.PaymentMethod)
}

func (s *ServiceSuite) TestGetOrderByUUID_WithCancelledOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 500),
		Status:          orderV1.OrderDtoStatusCANCELLED,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(expectedOrder, nil)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOrder, result)
	assert.Equal(s.T(), orderV1.OrderDtoStatusCANCELLED, result.Status)
	assert.Equal(s.T(), uuid.Nil, result.TransactionUUID)
}

func (s *ServiceSuite) TestGetOrderByUUID_WithComplexOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
	}

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(2000, 10000),
		Status:          orderV1.OrderDtoStatusPAID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP,
	}

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(expectedOrder, nil)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOrder, result)
	assert.Equal(s.T(), orderUUID, result.OrderUUID)
	assert.Equal(s.T(), userUUID, result.UserUUID)
	assert.Len(s.T(), result.PartUUIDs, 4)
	assert.Equal(s.T(), partUUIDs, result.PartUUIDs)
	assert.Equal(s.T(), orderV1.OrderDtoStatusPAID, result.Status)
	assert.Equal(s.T(), transactionUUID, result.TransactionUUID)
	assert.Equal(s.T(), orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP, result.PaymentMethod)
}

func (s *ServiceSuite) TestGetOrderByUUID_WithDifferentPaymentMethods() {
	// Arrange
	ctx := context.Background()
	testCases := []struct {
		name           string
		paymentMethod  orderV1.OrderDtoPaymentMethod
		expectedMethod orderV1.OrderDtoPaymentMethod
	}{
		{
			name:           "Credit Card",
			paymentMethod:  orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD,
			expectedMethod: orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD,
		},
		{
			name:           "SBP",
			paymentMethod:  orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP,
			expectedMethod: orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP,
		},
		{
			name:           "Investor Money",
			paymentMethod:  orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY,
			expectedMethod: orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY,
		},
		{
			name:           "Unspecified",
			paymentMethod:  orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
			expectedMethod: orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// Генерируем уникальные данные для каждого теста
			testOrderUUID := uuid.MustParse(gofakeit.UUID())
			testUserUUID := uuid.MustParse(gofakeit.UUID())
			testTransactionUUID := uuid.MustParse(gofakeit.UUID())
			testPartUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}
			totalPrice := gofakeit.Float32Range(100, 1000)

			expectedOrder := model.Order{
				OrderUUID:       testOrderUUID,
				UserUUID:        testUserUUID,
				PartUUIDs:       testPartUUIDs,
				TotalPrice:      totalPrice,
				Status:          orderV1.OrderDtoStatusPAID,
				TransactionUUID: testTransactionUUID,
				PaymentMethod:   tc.paymentMethod,
			}

			s.orderRepository.EXPECT().
				GetOrderByUUID(ctx, testOrderUUID).
				Return(expectedOrder, nil)

			// Act
			result, err := s.service.GetOrderByUUID(ctx, testOrderUUID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedOrder, result)
			assert.Equal(t, tc.expectedMethod, result.PaymentMethod)
		})
	}
}

func (s *ServiceSuite) TestGetOrderByUUID_WithEmptyPartsList() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{}, // Пустой список деталей
		TotalPrice:      float32(0.0),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(expectedOrder, nil)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOrder, result)
	assert.Empty(s.T(), result.PartUUIDs)
	assert.Equal(s.T(), float32(0.0), result.TotalPrice)
}

func (s *ServiceSuite) TestGetOrderByUUID_WithHighPriceOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(50000, 100000), // Высокая цена
		Status:          orderV1.OrderDtoStatusPAID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD,
	}

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(expectedOrder, nil)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOrder, result)
	assert.Greater(s.T(), result.TotalPrice, float32(50000))
	assert.Less(s.T(), result.TotalPrice, float32(100000))
}

func (s *ServiceSuite) TestGetOrderByUUID_WithZeroPriceOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())

	expectedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{uuid.MustParse(gofakeit.UUID())},
		TotalPrice:      float32(0.0), // Нулевая цена
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	s.orderRepository.EXPECT().
		GetOrderByUUID(ctx, orderUUID).
		Return(expectedOrder, nil)

	// Act
	result, err := s.service.GetOrderByUUID(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOrder, result)
	assert.Equal(s.T(), float32(0.0), result.TotalPrice)
}
