package part

import (
	"context"
	"errors"

	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
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
