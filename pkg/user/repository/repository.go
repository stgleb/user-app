package repository

import "user-app/pkg/user"

type reader interface {
	FindById(id string) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
	FindAll() ([]*user.User, error)
}

type writer interface {
	Update(user *user.User) error
	Store(user *user.User) (string, error)
}

type Repository interface {
	reader
	writer
}
