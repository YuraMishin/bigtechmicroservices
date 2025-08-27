package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/service/mocks"
)

func TestNewAPI_NilService(t *testing.T) {
	api, err := NewAPI(nil)
	assert.Nil(t, api)
	assert.Error(t, err)
}

func TestNewAPI_Success(t *testing.T) {
	mockSvc := mocks.NewPaymentService(t)
	api, err := NewAPI(mockSvc)
	assert.NoError(t, err)
	assert.NotNil(t, api)
}
