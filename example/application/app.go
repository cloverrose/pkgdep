package application

import (
	"github.com/cloverrose/pkgdep/example/service"
)

type Server struct {
	userService    *service.UserService
	addressService *service.AddressService
	// db             *infra.DB // ERROR
}

func (s *Server) UpdateName(ID int, name string) {
	s.userService.UpdateName(ID, name)
}

func (s *Server) UpdateAddress(ID string, zipCode, prefecture string) {
	s.addressService.UpdateAddress(ID, zipCode, prefecture)
}
