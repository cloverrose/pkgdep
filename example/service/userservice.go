package service

import (
	"github.com/cloverrose/pkgdep/example/domain/user"
	"github.com/cloverrose/pkgdep/example/infra"
)

type UserService struct {
	db *infra.DB
}

func (s *UserService) UpdateName(ID int, newName string) {
	var u *user.User
	u = s.db.GetUser(ID)
	u.Name = newName
	s.db.SaveUser(u)
}
