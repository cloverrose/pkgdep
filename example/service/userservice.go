package service

import (
	"github.com/cloverrose/pkgdep/example/infra"
)

type UserService struct {
	db *infra.DB
}

func (s *UserService) UpdateName(ID int, newName string) {
	user := s.db.GetUser(ID)
	user.Name = newName
	s.db.SaveUser(user)
}
