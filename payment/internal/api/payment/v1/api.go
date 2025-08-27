package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YuraMishin/bigtechmicroservices/payment/internal/service"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) (*api, error) {
	if paymentService == nil {
		return nil, status.Errorf(codes.InvalidArgument, "paymentService is nil")
	}

	return &api{
		paymentService: paymentService,
	}, nil
}
