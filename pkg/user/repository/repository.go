package repository

import (
	"context"
	"user-app/pkg/user"
)

// Read-only interface for user entity
type userReader interface {
	FindById(context.Context, string) (*user.User, error)
	FindByEmail(context.Context, string) (*user.User, error)
}

// Write-only interface for user entity
type userWriter interface {
	Update(context.Context, *user.User) error
	Store(context.Context, *user.User) (string, error)
}

// Read-only interface for token
type tokenReader interface {
	GetByEmail(context.Context, string) (*user.Token, error)
	GetByToken(context.Context, string) (*user.Token, error)
}

// Write-only interface for token
type tokenWriter interface {
	StoreToken(context.Context, *user.Token) error
	DisableToken(context.Context, string) error
}

type Repository interface {
	userReader
	userWriter
	tokenReader
	tokenWriter
}

type MemoryOpts struct{}

// MySQLOpts collects options for mysql repository
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
