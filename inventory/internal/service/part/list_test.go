package part

import (
	"context"
	"errors"

	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoPkg "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/part"
)

func (s *ServiceSuite) TestListParts_RepositoryError() {
	s.partRepository.EXPECT().ListParts(context.Background(), model.PartsFilter{}).Return(nil, errors.New("db down"))
	_, err := s.service.ListParts(context.Background(), model.PartsFilter{})
	require.Error(s.T(), err)
}

func (s *ServiceSuite) TestListParts_Success() {
	expected := []model.Part{{UUID: "id1"}, {UUID: "id2"}}
	s.partRepository.EXPECT().ListParts(context.Background(), model.PartsFilter{}).Return(expected, nil)
	got, err := s.service.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected, got)
}

// Repository tests (reuse this file, no new files created)
func (s *ServiceSuite) TestRepository_ListParts_NoFilter() {
	repo := repoPkg.NewRepository()
	parts, err := repo.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(parts), 0)
}

func (s *ServiceSuite) TestRepository_ListParts_FilterByUUID() {
	repo := repoPkg.NewRepository()
	all, err := repo.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(all), 0)

	firstID := all[0].UUID
	filtered, err := repo.ListParts(context.Background(), model.PartsFilter{UUIDs: []string{firstID}})
	require.NoError(s.T(), err)
	require.Len(s.T(), filtered, 1)
	require.Equal(s.T(), firstID, filtered[0].UUID)
}

func (s *ServiceSuite) TestRepository_GetPart_Found() {
	repo := repoPkg.NewRepository()
	all, err := repo.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(all), 0)

	got, err := repo.GetPart(context.Background(), all[0].UUID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), all[0].UUID, got.UUID)
}

func (s *ServiceSuite) TestRepository_GetPart_NotFound() {
	repo := repoPkg.NewRepository()
	_, err := repo.GetPart(context.Background(), "non-existent")
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrPartNotFound)
}
