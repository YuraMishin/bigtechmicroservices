package part

import (
	"context"
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPart_Success() {
	// Arrange
	ctx := context.Background()
	u := gofakeit.UUID()
	expectedPart := model.Part{
		UUID:          u,
		Name:          gofakeit.RandomString([]string{"Main Engine", "Auxiliary Engine", "Thruster", "Propulsion Unit"}),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(5000, 100000),
		StockQuantity: int64(gofakeit.IntRange(1, 20)),
		Category:      model.Category(gofakeit.IntRange(1, 4)),
		Dimensions: &model.Dimensions{
			Length: gofakeit.Float64Range(0.5, 10.0),
			Width:  gofakeit.Float64Range(0.3, 5.0),
			Height: gofakeit.Float64Range(0.1, 3.0),
			Weight: gofakeit.Float64Range(50, 2000),
		},
		Manufacturer: &model.Manufacturer{
			Name:    gofakeit.Company() + " Aerospace",
			Country: gofakeit.RandomString([]string{"USA", "Germany", "France", "Japan", "UK"}),
			Website: "https://" + gofakeit.DomainName(),
		},
		Tags:      []string{gofakeit.RandomString([]string{"propulsion", "thrust", "combustion", "high-performance", "main-engine"})},
		Metadata:  map[string]*model.Value{},
		CreatedAt: gofakeit.Int64(),
		UpdatedAt: gofakeit.Int64(),
	}

	s.partRepository.EXPECT().
		GetPart(ctx, expectedPart.UUID).
		Return(expectedPart, nil)

	// Act
	result, err := s.service.GetPart(ctx, expectedPart.UUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedPart, result)
}

func (s *ServiceSuite) TestGetPart_NotFound() {
	// Arrange
	ctx := context.Background()
	u := gofakeit.UUID()
	s.partRepository.EXPECT().
		GetPart(ctx, u).
		Return(model.Part{}, model.ErrPartNotFound)

	// Act
	result, err := s.service.GetPart(ctx, u)

	// Assert
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
	assert.Equal(s.T(), model.Part{}, result)
}

func (s *ServiceSuite) TestGetPart_RepositoryError() {
	// Arrange
	ctx := context.Background()
	u := gofakeit.UUID()
	expectedError := errors.New("database connection failed")
	s.partRepository.EXPECT().
		GetPart(ctx, u).
		Return(model.Part{}, expectedError)

	// Act
	result, err := s.service.GetPart(ctx, u)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Part{}, result)
}

func (s *ServiceSuite) TestGetPart_EmptyUUID_InvalidRequest() {
	// Arrange
	ctx := context.Background()
	// No repository expectation: service should validate and return early

	// Act
	result, err := s.service.GetPart(ctx, "")

	// Assert
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInvalidRequest)
	assert.Equal(s.T(), model.Part{}, result)
}

func (s *ServiceSuite) TestGetPart_WithNilOptionalFields() {
	// Arrange
	ctx := context.Background()
	u := gofakeit.UUID()
	expectedPart := model.Part{
		UUID:          u,
		Name:          gofakeit.RandomString([]string{"Test Part", "Sample Part", "Demo Part"}),
		Description:   gofakeit.Sentence(8),
		Price:         gofakeit.Float64Range(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 10)),
		Category:      model.CategoryFuel,
		Dimensions:    nil,
		Manufacturer:  nil,
		Tags:          []string{},
		Metadata:      map[string]*model.Value{},
		CreatedAt:     gofakeit.Int64(),
		UpdatedAt:     gofakeit.Int64(),
	}

	s.partRepository.EXPECT().
		GetPart(ctx, expectedPart.UUID).
		Return(expectedPart, nil)

	// Act
	result, err := s.service.GetPart(ctx, expectedPart.UUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedPart, result)
	assert.Nil(s.T(), result.Dimensions)
	assert.Nil(s.T(), result.Manufacturer)
}

func (s *ServiceSuite) TestGetPart_WithComplexMetadata() {
	// Arrange
	ctx := context.Background()
	u := gofakeit.UUID()
	thrustValue := gofakeit.Float64Range(500000, 2000000)
	fuelType := gofakeit.RandomString([]string{"liquid_hydrogen", "kerosene", "methane", "solid_fuel"})
	maxTemp := gofakeit.Float64Range(1500, 2500)
	isCertified := gofakeit.Bool()

	expectedPart := model.Part{
		UUID:          u,
		Name:          gofakeit.RandomString([]string{"Advanced Engine", "High-Performance Engine", "Turbine Engine"}),
		Description:   gofakeit.Sentence(12),
		Price:         gofakeit.Float64Range(50000, 200000),
		StockQuantity: int64(gofakeit.IntRange(1, 5)),
		Category:      model.CategoryEngine,
		Dimensions: &model.Dimensions{
			Length: gofakeit.Float64Range(10, 20),
			Width:  gofakeit.Float64Range(5, 10),
			Height: gofakeit.Float64Range(3, 6),
			Weight: gofakeit.Float64Range(500, 1500),
		},
		Manufacturer: &model.Manufacturer{
			Name:    gofakeit.Company() + " Aerospace",
			Country: gofakeit.RandomString([]string{"USA", "Germany", "France", "Japan"}),
			Website: "https://" + gofakeit.DomainName(),
		},
		Tags: []string{gofakeit.RandomString([]string{"high-performance", "turbine", "aviation"})},
		Metadata: map[string]*model.Value{
			"thrust": {
				DoubleValue: &thrustValue,
			},
			"fuel_type": {
				StringValue: &fuelType,
			},
			"max_temperature": {
				DoubleValue: &maxTemp,
			},
			"is_certified": {
				BoolValue: &isCertified,
			},
		},
		CreatedAt: gofakeit.Int64(),
		UpdatedAt: gofakeit.Int64(),
	}

	s.partRepository.EXPECT().
		GetPart(ctx, expectedPart.UUID).
		Return(expectedPart, nil)

	// Act
	result, err := s.service.GetPart(ctx, expectedPart.UUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedPart, result)
	assert.Len(s.T(), result.Metadata, 4)
	assert.NotNil(s.T(), result.Metadata["thrust"])
	assert.NotNil(s.T(), result.Metadata["fuel_type"])
	assert.NotNil(s.T(), result.Metadata["max_temperature"])
	assert.NotNil(s.T(), result.Metadata["is_certified"])
	assert.Equal(s.T(), thrustValue, *result.Metadata["thrust"].DoubleValue)
	assert.Equal(s.T(), fuelType, *result.Metadata["fuel_type"].StringValue)
	assert.Equal(s.T(), maxTemp, *result.Metadata["max_temperature"].DoubleValue)
	assert.Equal(s.T(), isCertified, *result.Metadata["is_certified"].BoolValue)
}
