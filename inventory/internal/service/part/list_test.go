package part

import (
	"context"
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
)

func (s *ServiceSuite) TestListParts_Success_EmptyFilter() {
	// Arrange
	ctx := context.Background()
	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Fuel Tank Alpha", "Fuel Tank Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(5000, 20000),
			StockQuantity: int64(gofakeit.IntRange(1, 15)),
			Category:      model.CategoryFuel,
			Tags:          []string{"storage"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	// Пустой фильтр - все поля nil или пустые
	filter := model.PartsFilter{}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 nil,
			Categories:            []repoModel.Category{},
			ManufacturerCountries: nil,
			Tags:                  nil,
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 2)
}

func (s *ServiceSuite) TestListParts_FilterByUUIDs() {
	// Arrange
	ctx := context.Background()
	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()

	expectedParts := []model.Part{
		{
			UUID:          uuid1,
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	filter := model.PartsFilter{
		UUIDs: []string{uuid1, uuid2},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 []string{uuid1, uuid2},
			Names:                 nil,
			Categories:            []repoModel.Category{},
			ManufacturerCountries: nil,
			Tags:                  nil,
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 1)
	assert.Equal(s.T(), uuid1, result[0].UUID)
}

func (s *ServiceSuite) TestListParts_FilterByNames() {
	// Arrange
	ctx := context.Background()
	name1 := "Engine Alpha"
	name2 := "Fuel Tank Beta"

	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          name1,
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
		{
			UUID:          gofakeit.UUID(),
			Name:          name2,
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(5000, 20000),
			StockQuantity: int64(gofakeit.IntRange(1, 15)),
			Category:      model.CategoryFuel,
			Tags:          []string{"storage"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	filter := model.PartsFilter{
		Names: []string{name1, name2},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 []string{name1, name2},
			Categories:            []repoModel.Category{},
			ManufacturerCountries: nil,
			Tags:                  nil,
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), name1, result[0].Name)
	assert.Equal(s.T(), name2, result[1].Name)
}

func (s *ServiceSuite) TestListParts_FilterByCategories() {
	// Arrange
	ctx := context.Background()
	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Gamma", "Engine Delta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(15000, 60000),
			StockQuantity: int64(gofakeit.IntRange(1, 8)),
			Category:      model.CategoryEngine,
			Tags:          []string{"thrust"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	filter := model.PartsFilter{
		Categories: []model.Category{model.CategoryEngine},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 nil,
			Categories:            []repoModel.Category{repoModel.CategoryEngine},
			ManufacturerCountries: nil,
			Tags:                  nil,
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), model.CategoryEngine, result[0].Category)
	assert.Equal(s.T(), model.CategoryEngine, result[1].Category)
}

func (s *ServiceSuite) TestListParts_FilterByManufacturerCountries() {
	// Arrange
	ctx := context.Background()
	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company() + " Aerospace",
				Country: "USA",
				Website: "https://" + gofakeit.DomainName(),
			},
			Tags:      []string{"propulsion"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: gofakeit.Int64(),
			UpdatedAt: gofakeit.Int64(),
		},
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Gamma", "Engine Delta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(15000, 60000),
			StockQuantity: int64(gofakeit.IntRange(1, 8)),
			Category:      model.CategoryEngine,
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company() + " Aerospace",
				Country: "Germany",
				Website: "https://" + gofakeit.DomainName(),
			},
			Tags:      []string{"thrust"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: gofakeit.Int64(),
			UpdatedAt: gofakeit.Int64(),
		},
	}

	filter := model.PartsFilter{
		ManufacturerCountries: []string{"USA", "Germany"},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 nil,
			Categories:            []repoModel.Category{},
			ManufacturerCountries: []string{"USA", "Germany"},
			Tags:                  nil,
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), "USA", result[0].Manufacturer.Country)
	assert.Equal(s.T(), "Germany", result[1].Manufacturer.Country)
}

func (s *ServiceSuite) TestListParts_FilterByTags() {
	// Arrange
	ctx := context.Background()
	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"thrust", "propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Gamma", "Engine Delta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(15000, 60000),
			StockQuantity: int64(gofakeit.IntRange(1, 8)),
			Category:      model.CategoryEngine,
			Tags:          []string{"thrust", "combustion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	filter := model.PartsFilter{
		Tags: []string{"thrust", "propulsion"},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 nil,
			Categories:            []repoModel.Category{},
			ManufacturerCountries: nil,
			Tags:                  []string{"thrust", "propulsion"},
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 2)
	assert.Contains(s.T(), result[0].Tags, "thrust")
	assert.Contains(s.T(), result[0].Tags, "propulsion")
	assert.Contains(s.T(), result[1].Tags, "thrust")
}

func (s *ServiceSuite) TestListParts_CombinedFilter() {
	// Arrange
	ctx := context.Background()
	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          "Engine Alpha",
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company() + " Aerospace",
				Country: "USA",
				Website: "https://" + gofakeit.DomainName(),
			},
			Tags:      []string{"thrust", "propulsion"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: gofakeit.Int64(),
			UpdatedAt: gofakeit.Int64(),
		},
	}

	filter := model.PartsFilter{
		Names:                 []string{"Engine Alpha"},
		Categories:            []model.Category{model.CategoryEngine},
		ManufacturerCountries: []string{"USA"},
		Tags:                  []string{"thrust"},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 []string{"Engine Alpha"},
			Categories:            []repoModel.Category{repoModel.CategoryEngine},
			ManufacturerCountries: []string{"USA"},
			Tags:                  []string{"thrust"},
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 1)
	assert.Equal(s.T(), "Engine Alpha", result[0].Name)
	assert.Equal(s.T(), model.CategoryEngine, result[0].Category)
	assert.Equal(s.T(), "USA", result[0].Manufacturer.Country)
	assert.Contains(s.T(), result[0].Tags, "thrust")
}

func (s *ServiceSuite) TestListParts_EmptyResult() {
	// Arrange
	ctx := context.Background()
	filter := model.PartsFilter{
		UUIDs: []string{"non-existent-uuid"},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 []string{"non-existent-uuid"},
			Names:                 nil,
			Categories:            []repoModel.Category{},
			ManufacturerCountries: nil,
			Tags:                  nil,
		}).
		Return([]model.Part{}, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), result)
}

func (s *ServiceSuite) TestListParts_RepositoryError() {
	// Arrange
	ctx := context.Background()
	expectedError := errors.New("database connection failed")
	filter := model.PartsFilter{
		Categories: []model.Category{model.CategoryEngine},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 nil,
			Categories:            []repoModel.Category{repoModel.CategoryEngine},
			ManufacturerCountries: nil,
			Tags:                  nil,
		}).
		Return(nil, expectedError)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), result)
}

func (s *ServiceSuite) TestListParts_EmptyFilter() {
	// Arrange
	ctx := context.Background()
	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	// Пустой фильтр - все поля nil или пустые
	filter := model.PartsFilter{}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 nil,
			Categories:            []repoModel.Category{},
			ManufacturerCountries: nil,
			Tags:                  nil,
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
}

func (s *ServiceSuite) TestListParts_FilterWithNilManufacturer() {
	// Arrange
	ctx := context.Background()
	expectedParts := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.RandomString([]string{"Engine Alpha", "Engine Beta"}),
			Description:   gofakeit.Sentence(8),
			Price:         gofakeit.Float64Range(10000, 50000),
			StockQuantity: int64(gofakeit.IntRange(1, 10)),
			Category:      model.CategoryEngine,
			Manufacturer:  nil, // nil manufacturer
			Tags:          []string{"propulsion"},
			Metadata:      map[string]*model.Value{},
			CreatedAt:     gofakeit.Int64(),
			UpdatedAt:     gofakeit.Int64(),
		},
	}

	filter := model.PartsFilter{
		ManufacturerCountries: []string{"USA"},
	}

	s.partRepository.EXPECT().
		ListParts(ctx, repoModel.PartsFilter{
			UUIDs:                 nil,
			Names:                 nil,
			Categories:            []repoModel.Category{},
			ManufacturerCountries: []string{"USA"},
			Tags:                  nil,
		}).
		Return(expectedParts, nil)

	// Act
	result, err := s.service.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedParts, result)
	assert.Len(s.T(), result, 1)
	assert.Nil(s.T(), result[0].Manufacturer)
}
