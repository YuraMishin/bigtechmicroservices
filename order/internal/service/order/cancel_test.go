package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestCancelOrderByUUID_Success_PendingPayment() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID()), uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       uuid.MustParse(orderUUID),
		UserUUID:        uuid.MustParse(userUUID),
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	s.orderRepository.EXPECT().
		UpdateOrder(ctx, model.Order{
			OrderUUID:       uuid.MustParse(orderUUID),
			UserUUID:        uuid.MustParse(userUUID),
			PartUUIDs:       partUUIDs,
			TotalPrice:      order.TotalPrice,
			Status:          orderV1.OrderDtoStatusCANCELLED, // Статус должен измениться на CANCELLED
			TransactionUUID: uuid.Nil,
			PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
		})

	// Act
	result, err := s.service.CancelOrderByUUID(ctx, order)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CancelOrderByUUIDNoContent{}, result)
}

func (s *ServiceSuite) TestCancelOrderByUUID_AlreadyPaid() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()

	order := model.Order{
		OrderUUID:       uuid.MustParse(orderUUID),
		UserUUID:        uuid.MustParse(userUUID),
		PartUUIDs:       []uuid.UUID{uuid.MustParse(gofakeit.UUID())},
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPAID,
		TransactionUUID: uuid.MustParse(transactionUUID),
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD,
	}

	// Act
	result, err := s.service.CancelOrderByUUID(ctx, order)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.Conflict{}, result)

	conflictError := result.(*orderV1.Conflict)
	assert.Equal(s.T(), 409, conflictError.Code)
	assert.Equal(s.T(), "Order is already paid and cannot be cancelled", conflictError.Message)
}

func (s *ServiceSuite) TestCancelOrderByUUID_AlreadyCancelled() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	order := model.Order{
		OrderUUID:       uuid.MustParse(orderUUID),
		UserUUID:        uuid.MustParse(userUUID),
		PartUUIDs:       []uuid.UUID{uuid.MustParse(gofakeit.UUID())},
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusCANCELLED,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	// Act
	result, err := s.service.CancelOrderByUUID(ctx, order)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CancelOrderByUUIDNoContent{}, result)
}

func (s *ServiceSuite) TestCancelOrderByUUID_InvalidStatus() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	order := model.Order{
		OrderUUID:       uuid.MustParse(orderUUID),
		UserUUID:        uuid.MustParse(userUUID),
		PartUUIDs:       []uuid.UUID{uuid.MustParse(gofakeit.UUID())},
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatus("INVALID_STATUS"), // Неизвестный статус
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	// Act
	result, err := s.service.CancelOrderByUUID(ctx, order)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.BadRequestError{}, result)

	badRequestError := result.(*orderV1.BadRequestError)
	assert.Equal(s.T(), 400, badRequestError.Code)
	assert.Equal(s.T(), "Invalid order status", badRequestError.Message)
}

func (s *ServiceSuite) TestCancelOrderByUUID_WithComplexOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUUIDs := []uuid.UUID{
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
	}

	order := model.Order{
		OrderUUID:       uuid.MustParse(orderUUID),
		UserUUID:        uuid.MustParse(userUUID),
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(1000, 5000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP,
	}

	// Ожидаем, что заказ будет обновлен со статусом CANCELLED
	expectedUpdatedOrder := model.Order{
		OrderUUID:       uuid.MustParse(orderUUID),
		UserUUID:        uuid.MustParse(userUUID),
		PartUUIDs:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		Status:          orderV1.OrderDtoStatusCANCELLED,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP,
	}

	s.orderRepository.EXPECT().
		UpdateOrder(ctx, expectedUpdatedOrder)

	// Act
	result, err := s.service.CancelOrderByUUID(ctx, order)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CancelOrderByUUIDNoContent{}, result)
}
