package service

import (
	"github.com/cloverrose/pkgdep/example/domain/address"
	"github.com/cloverrose/pkgdep/example/infra"
)

type AddressService struct {
	db *infra.DB
}

func (s *AddressService) UpdateAddress(ID string, zipCode, prefecture string) {
	var a *address.Address
	a = s.db.GetAddress(ID)
	a.ZipCode = zipCode
	a.Prefecture = prefecture
	s.db.SaveAddress(a)
}
