package repository

import (
	"context"
	"user-app/pkg/user"
)

type userReader interface {
	FindById(context.Context, string) (*user.User, error)
	FindByEmail(context.Context, string) (*user.User, error)
}

type userWriter interface {
	Update(context.Context, *user.User) error
	Store(context.Context, *user.User) (string, error)
}

type tokenRepository interface {
	StoreToken(context.Context, *user.Token) error
	GetByEmail(context.Context, string) (*user.Token, error)
	GetByToken(context.Context, string) (*user.Token, error)
	DisableToken(context.Context, string) error
}

type Repository interface {
	userReader
	userWriter
	tokenRepository
}

type MemoryOpts struct{}

type MySQLOpts struct {
	Protocol     string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type Opts struct {
	Mysql  MySQLOpts
	Memory MemoryOpts
}
