package part

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/mocks"
)

func TestNewService_NilRepository(t *testing.T) {
	// Act
	s, err := NewService(nil)

	// Assert
	assert.Nil(t, s)
	assert.Error(t, err)
}

func TestNewService_Success(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewPartRepository(t)

	// Act
	s, err := NewService(mockRepo)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, s)
}
