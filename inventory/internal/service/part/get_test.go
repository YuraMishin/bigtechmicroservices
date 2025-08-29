package part

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPart_InvalidUUID() {
	_, err := s.service.GetPart(context.Background(), "not-a-uuid")
	require.ErrorIs(s.T(), err, model.ErrInvalidRequest)
}

func (s *ServiceSuite) TestGetPart_RepositoryError() {
	id := uuid.New().String()
	s.partRepository.EXPECT().GetPart(context.Background(), id).Return(model.Part{}, errors.New("db down"))

	_, err := s.service.GetPart(context.Background(), id)
	require.Error(s.T(), err)
}

func (s *ServiceSuite) TestGetPart_Success() {
	id := uuid.New().String()
	expected := model.Part{UUID: id, Name: "warp coil"}
	s.partRepository.EXPECT().GetPart(context.Background(), id).Return(expected, nil)

	got, err := s.service.GetPart(context.Background(), id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected, got)
}
