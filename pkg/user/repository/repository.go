package repository

import "user-app/pkg/user"

type userReader interface {
	FindById(id string) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
	FindAll() ([]*user.User, error)
}

type userWriter interface {
	Update(user *user.User) error
	Store(user *user.User) (string, error)
}

type tokenRepository interface {
	StoreToken(*user.Token) error
	GetByEmail(string) (*user.Token, error)
	GetByToken(string) (*user.Token, error)
	RevokeToken(string) error
}

type Repository interface {
	userReader
	userWriter
	tokenRepository
}
