package service

import (
	"github.com/cloverrose/pkgdep/example/foundations/reservation/domain"
	"github.com/cloverrose/pkgdep/example/foundations/reservation/infra"
)

type ReservationService struct {
	d domain.ReservationDomain
	i infra.ReservationInfra
}
