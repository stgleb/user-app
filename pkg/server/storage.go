package server

import (
	"errors"
	"user-app/pkg/user/repository"
	"user-app/pkg/user/repository/memory"
	"user-app/pkg/user/repository/mysql"
)

const (
	MySQL    = "mysql"
	InMemory = "memory"
)

func NewRepository(repositoryType string, opts repository.Opts) (repository.Repository, error) {
	switch repositoryType {
	case MySQL:
		r, err := mysql.NewRepository(opts.Mysql)
		if err != nil {
			return nil, err
		}
		return r, nil
	case InMemory:
		return memory.NewRepository(), nil
	}
	return nil, errors.New("unknown repository type")
}
