package service

import (
	"github.com/cloverrose/pkgdep/example/infra"
)

type AddressService struct {
	db *infra.DB
}

func (s *AddressService) UpdateAddress(ID string, zipCode, prefecture string) {
	a := s.db.GetAddress(ID)
	a.ZipCode = zipCode
	a.Prefecture = prefecture
	s.db.SaveAddress(a)
}
