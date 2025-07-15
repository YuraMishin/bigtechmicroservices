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

func (s *ServiceSuite) TestCreateNewOrder_Success() {
	// Arrange
	ctx := context.Background()
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUID1 := uuid.MustParse(gofakeit.UUID())
	partUUID2 := uuid.MustParse(gofakeit.UUID())

	request := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID1, partUUID2},
	}

	// Мокируем детали из inventory service
	inventoryParts := []model.Part{
		{
			UUID:          partUUID1.String(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(1000, 5000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
		{
			UUID:          partUUID2.String(),
			Name:          gofakeit.RandomString([]string{"Fuel Tank Alpha", "Fuel Tank Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(500, 2000),
			StockQuantity: int64(gofakeit.IntRange(1, 15)),
			Category:      model.CategoryFuel,
			Tags:          []string{"storage"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	expectedTotalPrice := float32(inventoryParts[0].Price + inventoryParts[1].Price)

	// Ожидаем вызов inventory client
	s.inventoryClient.EXPECT().
		ListParts(ctx, model.PartsFilter{
			UUIDs: []string{partUUID1.String(), partUUID2.String()},
		}).
		Return(inventoryParts, nil)

	// Ожидаем создание заказа в репозитории
	s.orderRepository.EXPECT().
		CreateNewOrder(ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.UserUUID == userUUID &&
				len(order.PartUUIDs) == 2 &&
				order.PartUUIDs[0] == partUUID1 &&
				order.PartUUIDs[1] == partUUID2 &&
				order.TotalPrice > 0 && // Проверяем, что цена положительная
				order.Status == orderV1.OrderDtoStatusPENDINGPAYMENT &&
				order.TransactionUUID == uuid.Nil &&
				order.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
		}))

	// Act
	result, err := s.service.CreateNewOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CreateOrderResponse{}, result)

	createResponse := result.(*orderV1.CreateOrderResponse)
	assert.NotNil(s.T(), createResponse.OrderUUID)
	assert.InDelta(s.T(), float64(expectedTotalPrice), float64(createResponse.TotalPrice), 0.01)
}

func (s *ServiceSuite) TestCreateNewOrder_SinglePart() {
	// Arrange
	ctx := context.Background()
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUID := uuid.MustParse(gofakeit.UUID())

	request := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID},
	}

	inventoryPart := model.Part{
		UUID:          partUUID.String(),
		Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
		Description:   gofakeit.Sentence(8),
		Price:         gofakeit.Float64Range(1000, 5000),
		StockQuantity: int64(gofakeit.IntRange(1, 10)),
		Category:      model.CategoryEngine,
		Tags:          []string{"propulsion"},
		Metadata:      map[string]*model.Value{},
		CreatedAt:     gofakeit.Int64(),
		UpdatedAt:     gofakeit.Int64(),
	}

	expectedTotalPrice := float32(inventoryPart.Price)

	s.inventoryClient.EXPECT().
		ListParts(ctx, model.PartsFilter{
			UUIDs: []string{partUUID.String()},
		}).
		Return([]model.Part{inventoryPart}, nil)

	s.orderRepository.EXPECT().
		CreateNewOrder(ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.UserUUID == userUUID &&
				len(order.PartUUIDs) == 1 &&
				order.PartUUIDs[0] == partUUID &&
				order.TotalPrice > 0 && // Проверяем, что цена положительная
				order.Status == orderV1.OrderDtoStatusPENDINGPAYMENT &&
				order.TransactionUUID == uuid.Nil &&
				order.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
		}))

	// Act
	result, err := s.service.CreateNewOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CreateOrderResponse{}, result)

	createResponse := result.(*orderV1.CreateOrderResponse)
	assert.NotNil(s.T(), createResponse.OrderUUID)
	assert.InDelta(s.T(), float64(expectedTotalPrice), float64(createResponse.TotalPrice), 0.01)
}

func (s *ServiceSuite) TestCreateNewOrder_PartsNotFound() {
	// Arrange
	ctx := context.Background()
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUID1 := uuid.MustParse(gofakeit.UUID())
	partUUID2 := uuid.MustParse(gofakeit.UUID())

	request := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID1, partUUID2},
	}

	// Inventory возвращает только одну деталь вместо двух
	inventoryParts := []model.Part{
		{
			UUID:          partUUID1.String(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(1000, 5000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	s.inventoryClient.EXPECT().
		ListParts(ctx, model.PartsFilter{
			UUIDs: []string{partUUID1.String(), partUUID2.String()},
		}).
		Return(inventoryParts, nil)

	// Act
	result, err := s.service.CreateNewOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.BadRequestError{}, result)

	badRequestError := result.(*orderV1.BadRequestError)
	assert.Equal(s.T(), 400, badRequestError.Code)
	assert.Equal(s.T(), "Some parts not found in inventory", badRequestError.Message)
}

func (s *ServiceSuite) TestCreateNewOrder_InventoryClientError() {
	// Arrange
	ctx := context.Background()
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUID := uuid.MustParse(gofakeit.UUID())

	request := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID},
	}

	s.inventoryClient.EXPECT().
		ListParts(ctx, model.PartsFilter{
			UUIDs: []string{partUUID.String()},
		}).
		Return(nil, assert.AnError)

	// Act
	result, err := s.service.CreateNewOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.InternalServerError{}, result)

	internalError := result.(*orderV1.InternalServerError)
	assert.Equal(s.T(), 500, internalError.Code)
	assert.Equal(s.T(), "Internal error", internalError.Message)
}

func (s *ServiceSuite) TestCreateNewOrder_EmptyPartsList() {
	// Arrange
	ctx := context.Background()
	userUUID := uuid.MustParse(gofakeit.UUID())

	request := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{},
	}

	s.inventoryClient.EXPECT().
		ListParts(ctx, model.PartsFilter{
			UUIDs: []string{},
		}).
		Return([]model.Part{}, nil)

	// Ожидаем создание заказа в репозитории
	s.orderRepository.EXPECT().
		CreateNewOrder(ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.UserUUID == userUUID &&
				len(order.PartUUIDs) == 0 &&
				order.TotalPrice == float32(0.0) &&
				order.Status == orderV1.OrderDtoStatusPENDINGPAYMENT &&
				order.TransactionUUID == uuid.Nil &&
				order.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
		}))

	// Act
	result, err := s.service.CreateNewOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CreateOrderResponse{}, result)

	createResponse := result.(*orderV1.CreateOrderResponse)
	assert.NotNil(s.T(), createResponse.OrderUUID)
	assert.Equal(s.T(), float32(0.0), createResponse.TotalPrice)
}

func (s *ServiceSuite) TestCreateNewOrder_WithComplexParts() {
	// Arrange
	ctx := context.Background()
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUIDs := []uuid.UUID{
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
	}

	request := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	}

	// Создаем сложные детали с различными категориями и ценами
	inventoryParts := []model.Part{
		{
			UUID:          partUUIDs[0].String(),
			Name:          gofakeit.RandomString([]string{"Main Engine", "Auxiliary Engine"}),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(5000, 15000),
			StockQuantity: int64(gofakeit.IntRange(1, 5)),
			Category:      model.CategoryEngine,
			Dimensions: &model.Dimensions{
				Length: gofakeit.Float64Range(1, 5),
				Width:  gofakeit.Float64Range(0.5, 2),
				Height: gofakeit.Float64Range(0.5, 2),
				Weight: gofakeit.Float64Range(100, 500),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company() + " Aerospace",
				Country: gofakeit.RandomString([]string{"USA", "Germany"}),
				Website: "https://" + gofakeit.DomainName(),
			},
			Tags:      []string{"propulsion", "thrust", "high-performance"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: gofakeit.Int64(),
			UpdatedAt: gofakeit.Int64(),
		},
		{
			UUID:          partUUIDs[1].String(),
			Name:          gofakeit.RandomString([]string{"Fuel Tank", "Cryogenic Tank"}),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(2000, 8000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryFuel,
			Dimensions: &model.Dimensions{
				Length: gofakeit.Float64Range(2, 8),
				Width:  gofakeit.Float64Range(1, 3),
				Height: gofakeit.Float64Range(1, 3),
				Weight: gofakeit.Float64Range(200, 800),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company() + " Systems",
				Country: gofakeit.RandomString([]string{"France", "Japan"}),
				Website: "https://" + gofakeit.DomainName(),
			},
			Tags:      []string{"storage", "cryogenic", "fuel"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: gofakeit.Int64(),
			UpdatedAt: gofakeit.Int64(),
		},
		{
			UUID:          partUUIDs[2].String(),
			Name:          gofakeit.RandomString([]string{"Main Wing", "Control Surface"}),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(3000, 12000),
			StockQuantity: int64(gofakeit.IntRange(1, 8)),
			Category:      model.CategoryWing,
			Dimensions: &model.Dimensions{
				Length: gofakeit.Float64Range(5, 15),
				Width:  gofakeit.Float64Range(0.1, 1),
				Height: gofakeit.Float64Range(0.1, 1),
				Weight: gofakeit.Float64Range(50, 300),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company() + " Dynamics",
				Country: gofakeit.RandomString([]string{"UK", "Canada"}),
				Website: "https://" + gofakeit.DomainName(),
			},
			Tags:      []string{"aerodynamics", "lift", "control"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: gofakeit.Int64(),
			UpdatedAt: gofakeit.Int64(),
		},
	}

	expectedTotalPrice := float32(inventoryParts[0].Price + inventoryParts[1].Price + inventoryParts[2].Price)

	s.inventoryClient.EXPECT().
		ListParts(ctx, model.PartsFilter{
			UUIDs: []string{partUUIDs[0].String(), partUUIDs[1].String(), partUUIDs[2].String()},
		}).
		Return(inventoryParts, nil)

	s.orderRepository.EXPECT().
		CreateNewOrder(ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.UserUUID == userUUID &&
				len(order.PartUUIDs) == 3 &&
				order.PartUUIDs[0] == partUUIDs[0] &&
				order.PartUUIDs[1] == partUUIDs[1] &&
				order.PartUUIDs[2] == partUUIDs[2] &&
				order.TotalPrice > 0 && // Проверяем, что цена положительная
				order.Status == orderV1.OrderDtoStatusPENDINGPAYMENT &&
				order.TransactionUUID == uuid.Nil &&
				order.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
		}))

	// Act
	result, err := s.service.CreateNewOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CreateOrderResponse{}, result)

	createResponse := result.(*orderV1.CreateOrderResponse)
	assert.NotNil(s.T(), createResponse.OrderUUID)
	assert.InDelta(s.T(), float64(expectedTotalPrice), float64(createResponse.TotalPrice), 0.01)
}

func (s *ServiceSuite) TestCreateNewOrder_WithNilOptionalFields() {
	// Arrange
	ctx := context.Background()
	userUUID := uuid.MustParse(gofakeit.UUID())
	partUUID := uuid.MustParse(gofakeit.UUID())

	request := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID},
	}

	// Деталь с nil опциональными полями
	inventoryPart := model.Part{
		UUID:          partUUID.String(),
		Name:          gofakeit.RandomString([]string{"Simple Part"}),
		Description:   gofakeit.Sentence(8),
		Price:         gofakeit.Float64Range(100, 500),
		StockQuantity: int64(gofakeit.IntRange(1, 10)),
		Category:      model.CategoryPorthole,
		Dimensions:    nil, // nil optional field
		Manufacturer:  nil, // nil optional field
		Tags:          []string{},
		Metadata:      map[string]*model.Value{},
		CreatedAt:     gofakeit.Int64(),
		UpdatedAt:     gofakeit.Int64(),
	}

	expectedTotalPrice := float32(inventoryPart.Price)

	s.inventoryClient.EXPECT().
		ListParts(ctx, model.PartsFilter{
			UUIDs: []string{partUUID.String()},
		}).
		Return([]model.Part{inventoryPart}, nil)

	s.orderRepository.EXPECT().
		CreateNewOrder(ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.UserUUID == userUUID &&
				len(order.PartUUIDs) == 1 &&
				order.PartUUIDs[0] == partUUID &&
				order.TotalPrice > 0 && // Проверяем, что цена положительная
				order.Status == orderV1.OrderDtoStatusPENDINGPAYMENT &&
				order.TransactionUUID == uuid.Nil &&
				order.PaymentMethod == orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
		}))

	// Act
	result, err := s.service.CreateNewOrder(ctx, request)

	// Assert
	assert.NoError(s.T(), err)
	assert.IsType(s.T(), &orderV1.CreateOrderResponse{}, result)

	createResponse := result.(*orderV1.CreateOrderResponse)
	assert.NotNil(s.T(), createResponse.OrderUUID)
	assert.InDelta(s.T(), float64(expectedTotalPrice), float64(createResponse.TotalPrice), 0.01)
}
