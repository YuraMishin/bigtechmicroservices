package part

import (
	"context"
	"errors"

	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
)

func (s *ServiceSuite) TestListParts_Success() {
	// Arrange
	ctx := context.Background()
	filter := model.PartsFilter{
		UUIDs:                 []string{"id-1", "id-2"},
		Names:                 []string{"Engine", "Thruster"},
		Categories:            []model.Category{model.CategoryEngine, model.CategoryPorthole},
		ManufacturerCountries: []string{"USA"},
		Tags:                  []string{"high-performance"},
	}

	expected := []model.Part{{UUID: "id-1"}, {UUID: "id-2"}}

	s.partRepository.EXPECT().
		ListParts(ctx, filter).
		Return(expected, nil)

	// Act
	parts, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, parts)
}

func (s *ServiceSuite) TestListParts_RepositoryError() {
	// Arrange
	ctx := context.Background()
	filter := model.PartsFilter{}
	expectedErr := errors.New("repo failed")

	s.partRepository.EXPECT().
		ListParts(ctx, filter).
		Return(nil, expectedErr)

	// Act
	parts, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedErr, err)
	assert.Nil(s.T(), parts)
}
