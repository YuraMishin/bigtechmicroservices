package app

import (
	"context"
	"fmt"

	paymentV1API "github.com/YuraMishin/bigtechmicroservices/payment/internal/api/payment/v1"
	"github.com/YuraMishin/bigtechmicroservices/payment/internal/service"
	paymentService "github.com/YuraMishin/bigtechmicroservices/payment/internal/service/payment"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	var err error
	if d.paymentV1API == nil {
		d.paymentV1API, err = paymentV1API.NewAPI(d.PaymentService(ctx))
		if err != nil {
			panic(fmt.Sprintf("failed to create paymentV1API: %v\n", err))
		}
	}
	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
