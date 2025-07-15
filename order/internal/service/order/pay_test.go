package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestPayOrder_Success_Card() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD,
	}

	paymentResponse := &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}

	// Ожидаем вызов payment client
	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_CARD").
		Return(paymentResponse, nil)

	// Ожидаем обновление заказа в репозитории
	s.orderRepository.EXPECT().
		UpdateOrder(ctx, mock.MatchedBy(func(updatedOrder model.Order) bool {
			return updatedOrder.OrderUUID == orderUUID &&
				updatedOrder.UserUUID == userUUID &&
				updatedOrder.PartUUIDs[0] == partUUIDs[0] &&
				updatedOrder.TotalPrice == order.TotalPrice &&
				updatedOrder.Status == orderV1.OrderDtoStatusPAID &&
				updatedOrder.TransactionUUID == transactionUUID &&
				updatedOrder.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD
		}))

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.PayOrderResponse{}, result)

	payResponse := result.(*orderV1.PayOrderResponse)
	assert.Equal(s.T(), transactionUUID, payResponse.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_Success_SBP() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP,
	}

	paymentResponse := &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}

	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_SBP").
		Return(paymentResponse, nil)

	s.orderRepository.EXPECT().
		UpdateOrder(ctx, mock.MatchedBy(func(updatedOrder model.Order) bool {
			return updatedOrder.OrderUUID == orderUUID &&
				updatedOrder.UserUUID == userUUID &&
				updatedOrder.PartUUIDs[0] == partUUIDs[0] &&
				updatedOrder.TotalPrice == order.TotalPrice &&
				updatedOrder.Status == orderV1.OrderDtoStatusPAID &&
				updatedOrder.TransactionUUID == transactionUUID &&
				updatedOrder.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP
		}))

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.PayOrderResponse{}, result)

	payResponse := result.(*orderV1.PayOrderResponse)
	assert.Equal(s.T(), transactionUUID, payResponse.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_Success_CreditCard() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD,
	}

	paymentResponse := &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}

	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_CREDIT_CARD").
		Return(paymentResponse, nil)

	s.orderRepository.EXPECT().
		UpdateOrder(ctx, mock.MatchedBy(func(updatedOrder model.Order) bool {
			return updatedOrder.OrderUUID == orderUUID &&
				updatedOrder.UserUUID == userUUID &&
				updatedOrder.PartUUIDs[0] == partUUIDs[0] &&
				updatedOrder.TotalPrice == order.TotalPrice &&
				updatedOrder.Status == orderV1.OrderDtoStatusPAID &&
				updatedOrder.TransactionUUID == transactionUUID &&
				updatedOrder.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD
		}))

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.PayOrderResponse{}, result)

	payResponse := result.(*orderV1.PayOrderResponse)
	assert.Equal(s.T(), transactionUUID, payResponse.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_Success_InvestorMoney() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY,
	}

	paymentResponse := &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}

	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_INVESTOR_MONEY").
		Return(paymentResponse, nil)

	s.orderRepository.EXPECT().
		UpdateOrder(ctx, mock.MatchedBy(func(updatedOrder model.Order) bool {
			return updatedOrder.OrderUUID == orderUUID &&
				updatedOrder.UserUUID == userUUID &&
				updatedOrder.PartUUIDs[0] == partUUIDs[0] &&
				updatedOrder.TotalPrice == order.TotalPrice &&
				updatedOrder.Status == orderV1.OrderDtoStatusPAID &&
				updatedOrder.TransactionUUID == transactionUUID &&
				updatedOrder.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY
		}))

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.PayOrderResponse{}, result)

	payResponse := result.(*orderV1.PayOrderResponse)
	assert.Equal(s.T(), transactionUUID, payResponse.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_UnspecifiedPaymentMethod() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.BadRequestError{}, result)

	badRequestError := result.(*orderV1.BadRequestError)
	assert.Equal(s.T(), 400, badRequestError.Code)
	assert.Equal(s.T(), "Payment method must be specified", badRequestError.Message)
}

func (s *ServiceSuite) TestPayOrder_UnknownPaymentMethod() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethod("UNKNOWN_METHOD"), // Неизвестный метод
	}

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.BadRequestError{}, result)

	badRequestError := result.(*orderV1.BadRequestError)
	assert.Equal(s.T(), 400, badRequestError.Code)
	assert.Equal(s.T(), "Unknown payment method", badRequestError.Message)
}

func (s *ServiceSuite) TestPayOrder_PaymentClientError() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD,
	}

	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_CARD").
		Return(nil, assert.AnError)

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.InternalServerError{}, result)

	internalError := result.(*orderV1.InternalServerError)
	assert.Equal(s.T(), 500, internalError.Code)
	assert.Equal(s.T(), "Internal error", internalError.Message)
}

func (s *ServiceSuite) TestPayOrder_UnexpectedResponseType() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(100, 1000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD,
	}

	// Возвращаем неожиданный тип ответа
	unexpectedResponse := &orderV1.BadRequestError{
		Code:    400,
		Message: "Unexpected response",
	}

	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_CARD").
		Return(unexpectedResponse, nil)

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.InternalServerError{}, result)

	internalError := result.(*orderV1.InternalServerError)
	assert.Equal(s.T(), 500, internalError.Code)
	assert.Equal(s.T(), "Internal error", internalError.Message)
}

func (s *ServiceSuite) TestPayOrder_WithComplexOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
	}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(2000, 10000),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP,
	}

	paymentResponse := &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}

	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_SBP").
		Return(paymentResponse, nil)

	s.orderRepository.EXPECT().
		UpdateOrder(ctx, mock.MatchedBy(func(updatedOrder model.Order) bool {
			return updatedOrder.OrderUUID == orderUUID &&
				updatedOrder.UserUUID == userUUID &&
				len(updatedOrder.PartUUIDs) == 3 &&
				updatedOrder.PartUUIDs[0] == partUUIDs[0] &&
				updatedOrder.PartUUIDs[1] == partUUIDs[1] &&
				updatedOrder.PartUUIDs[2] == partUUIDs[2] &&
				updatedOrder.TotalPrice == order.TotalPrice &&
				updatedOrder.Status == orderV1.OrderDtoStatusPAID &&
				updatedOrder.TransactionUUID == transactionUUID &&
				updatedOrder.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP
		}))

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.PayOrderResponse{}, result)

	payResponse := result.(*orderV1.PayOrderResponse)
	assert.Equal(s.T(), transactionUUID, payResponse.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_WithHighValueOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.MustParse(gofakeit.UUID())
	userUUID := uuid.MustParse(gofakeit.UUID())
	transactionUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{uuid.MustParse(gofakeit.UUID())}

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      gofakeit.Float32Range(50000, 100000), // Высокая стоимость
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
	}

	request := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD,
	}

	paymentResponse := &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}

	s.paymentClient.EXPECT().
		PayOrder(ctx, orderUUID.String(), userUUID.String(), "PAYMENT_METHOD_CREDIT_CARD").
		Return(paymentResponse, nil)

	s.orderRepository.EXPECT().
		UpdateOrder(ctx, mock.MatchedBy(func(updatedOrder model.Order) bool {
			return updatedOrder.OrderUUID == orderUUID &&
				updatedOrder.UserUUID == userUUID &&
				updatedOrder.PartUUIDs[0] == partUUIDs[0] &&
				updatedOrder.TotalPrice == order.TotalPrice &&
				updatedOrder.Status == orderV1.OrderDtoStatusPAID &&
				updatedOrder.TransactionUUID == transactionUUID &&
				updatedOrder.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD
		}))

	// Act
	result, err := s.service.PayOrder(ctx, order, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.PayOrderResponse{}, result)

	payResponse := result.(*orderV1.PayOrderResponse)
	assert.Equal(s.T(), transactionUUID, payResponse.TransactionUUID)
}
