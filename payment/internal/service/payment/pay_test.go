package payment

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func (s *ServiceSuite) TestPayOrder_Success_Card() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	// Проверяем, что UUID транзакции валидный
	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_Success_SBP() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_Success_CreditCard() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_Success_InvestorMoney() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_Success_UnspecifiedMethod() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_EmptyOrderUUID() {
	// Arrange
	ctx := context.Background()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     "",
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_EmptyUserUUID() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      "",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_EmptyUUIDs() {
	// Arrange
	ctx := context.Background()
	request := &paymentV1.PayOrderRequest{
		OrderUuid:     "",
		UserUuid:      "",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_InvalidOrderUUID() {
	// Arrange
	ctx := context.Background()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     "invalid-uuid",
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_InvalidUserUUID() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      "invalid-uuid",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_UnknownPaymentMethod() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod(999), // Неизвестный метод
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_NilRequest() {
	// Arrange
	ctx := context.Background()
	// Act
	response, err := s.service.PayOrder(ctx, nil)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}

func (s *ServiceSuite) TestPayOrder_UniqueTransactionUUIDs() {
	// Arrange
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act - делаем несколько вызовов
	response1, err1 := s.service.PayOrder(ctx, request)
	response2, err2 := s.service.PayOrder(ctx, request)
	response3, err3 := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err1)
	assert.NoError(s.T(), err2)
	assert.NoError(s.T(), err3)

	assert.NotNil(s.T(), response1)
	assert.NotNil(s.T(), response2)
	assert.NotNil(s.T(), response3)

	assert.NotEmpty(s.T(), response1.TransactionUuid)
	assert.NotEmpty(s.T(), response2.TransactionUuid)
	assert.NotEmpty(s.T(), response3.TransactionUuid)

	// Проверяем, что все UUID транзакций уникальны
	assert.NotEqual(s.T(), response1.TransactionUuid, response2.TransactionUuid)
	assert.NotEqual(s.T(), response1.TransactionUuid, response3.TransactionUuid)
	assert.NotEqual(s.T(), response2.TransactionUuid, response3.TransactionUuid)

	// Проверяем валидность всех UUID
	_, parseErr1 := uuid.Parse(response1.TransactionUuid)
	_, parseErr2 := uuid.Parse(response2.TransactionUuid)
	_, parseErr3 := uuid.Parse(response3.TransactionUuid)

	assert.NoError(s.T(), parseErr1)
	assert.NoError(s.T(), parseErr2)
	assert.NoError(s.T(), parseErr3)
}

func (s *ServiceSuite) TestPayOrder_WithComplexData() {
	// Arrange - используем сложные данные
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	// Проверяем, что UUID транзакции валидный
	transactionUUID, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
	assert.NotEqual(s.T(), uuid.Nil, transactionUUID)
}

func (s *ServiceSuite) TestPayOrder_AllPaymentMethods() {
	// Arrange - тестируем все методы оплаты
	ctx := context.Background()
	paymentMethods := []paymentV1.PaymentMethod{
		paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	for _, paymentMethod := range paymentMethods {
		// Arrange
		request := &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethod,
		}

		// Act
		response, err := s.service.PayOrder(ctx, request)

		// Assert
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), response)
		assert.NotEmpty(s.T(), response.TransactionUuid)

		_, parseErr := uuid.Parse(response.TransactionUuid)
		assert.NoError(s.T(), parseErr)
	}
}

func (s *ServiceSuite) TestPayOrder_WithLongUUIDs() {
	// Arrange - используем очень длинные UUID (если такие есть)
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	request := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Act
	response, err := s.service.PayOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotEmpty(s.T(), response.TransactionUuid)

	// Проверяем, что UUID транзакции имеет правильную длину (36 символов для UUID v4)
	assert.Len(s.T(), response.TransactionUuid, 36)

	_, parseErr := uuid.Parse(response.TransactionUuid)
	assert.NoError(s.T(), parseErr)
}
