package payment

import (
	def "github.com/YuraMishin/bigtechmicroservices/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
