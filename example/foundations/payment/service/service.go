package service

import (
	"github.com/cloverrose/pkgdep/example/foundations/payment/domain"
	"github.com/cloverrose/pkgdep/example/foundations/payment/infra"
)

type PaymentService struct {
	d domain.PaymentDomain
	i infra.PaymentInfra
}
