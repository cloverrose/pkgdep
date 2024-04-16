package infra

import (
	"github.com/cloverrose/pkgdep/example/domain/address"
	"github.com/cloverrose/pkgdep/example/domain/user"
)

type DB struct{}

func (d *DB) GetUser(ID int) *user.User {
	return &user.User{ID, "dummy-name", "dummy-address-ID"}
}

func (d *DB) SaveUser(user *user.User) {
	// Save user
}

func (d *DB) GetAddress(ID string) *address.Address {
	return &address.Address{ID, "123-4567", "Tokyo"}
}

func (d *DB) SaveAddress(address *address.Address) {
	// Save address
}
