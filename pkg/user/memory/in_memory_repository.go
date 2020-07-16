package memory

import (
	"errors"
	"user-app/pkg/user"
)

type InMemoryRepository struct {
	db map[string]*user.User
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		db: make(map[string]*user.User),
	}
}

func (r *InMemoryRepository) FindById(id string) (*user.User, error) {
	if _, ok := r.db[id]; !ok {
		return nil, user.NotFound
	}
	return r.db[id], nil
}

func (r *InMemoryRepository) FindByEmail(email string) (*user.User, error) {
	for _, u := range r.db {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, user.NotFound
}

func (r *InMemoryRepository) FindAll() ([]*user.User, error) {
	users := make([]*user.User, 0)
	for _, u := range r.db {
		users = append(users, u)
	}

	return users, nil
}

func (r *InMemoryRepository) Update(user *user.User) error {
	_, ok := r.db[user.Id]
	if !ok {
		return errors.New("not found")
	}

	r.db[user.Id] = user
	return nil
}

func (r *InMemoryRepository) Store(user *user.User) (string, error) {
	r.db[user.Id] = user
	return user.Id, nil
}

func init() {
	var _ user.Repository = &InMemoryRepository{}
}
