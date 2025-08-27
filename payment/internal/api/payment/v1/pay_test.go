package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/service/mocks"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

func TestAPI_PayOrder_Success(t *testing.T) {
	mockSvc := mocks.NewPaymentService(t)
	api, err := NewAPI(mockSvc)
	assert.NoError(t, err)

	req := &paymentV1.PayOrderRequest{OrderUuid: "o", UserUuid: "u"}
	mockSvc.EXPECT().PayOrder(context.Background(), req).Return(&paymentV1.PayOrderResponse{TransactionUuid: "txid"}, nil)

	resp, err := api.PayOrder(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "txid", resp.GetTransactionUuid())
}

func TestAPI_PayOrder_Error(t *testing.T) {
	mockSvc := mocks.NewPaymentService(t)
	api, err := NewAPI(mockSvc)
	assert.NoError(t, err)

	req := &paymentV1.PayOrderRequest{OrderUuid: "o", UserUuid: "u"}
	mockSvc.EXPECT().PayOrder(context.Background(), req).Return((*paymentV1.PayOrderResponse)(nil), errors.New("db down"))

	resp, err := api.PayOrder(context.Background(), req)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	if assert.True(t, ok) {
		assert.Equal(t, codes.Internal, st.Code())
	}
}
