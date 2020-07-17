package memory

import (
	"errors"
	"user-app/pkg/user"
	"user-app/pkg/user/repository"
)

type InMemoryRepository struct {
	userDb  map[string]*user.User
	tokenDb map[string]*user.Token
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		userDb:  make(map[string]*user.User),
		tokenDb: make(map[string]*user.Token),
	}
}

func (r *InMemoryRepository) FindById(id string) (*user.User, error) {
	if _, ok := r.userDb[id]; !ok {
		return nil, user.NotFound
	}
	return r.userDb[id], nil
}

func (r *InMemoryRepository) FindByEmail(email string) (*user.User, error) {
	for _, u := range r.userDb {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, user.NotFound
}

func (r *InMemoryRepository) FindAll() ([]*user.User, error) {
	users := make([]*user.User, 0)
	for _, u := range r.userDb {
		users = append(users, u)
	}

	return users, nil
}

func (r *InMemoryRepository) Update(user *user.User) error {
	_, ok := r.userDb[user.Id]
	if !ok {
		return errors.New("not found")
	}

	r.userDb[user.Id] = user
	return nil
}

func (r *InMemoryRepository) Store(user *user.User) (string, error) {
	r.userDb[user.Id] = user
	return user.Id, nil
}

func (r *InMemoryRepository) StoreToken(token *user.Token) error {
	r.tokenDb[token.Email] = token
	return nil
}

func (r *InMemoryRepository) GetByEmail(email string) (*user.Token, error) {
	token, ok := r.tokenDb[email]
	if !ok {
		return nil, user.NotFound
	}
	return token, nil
}

func (r *InMemoryRepository) GetByToken(tokenValue string) (*user.Token, error) {
	for _, token := range r.tokenDb {
		if token.Value == tokenValue {
			return token, nil
		}
	}
	return nil, user.NotFound
}

func (r *InMemoryRepository) RevokeToken(email string) error {
	t := r.tokenDb[email]
	t.Used = true
	return nil
}

func init() {
	var _ repository.Repository = &InMemoryRepository{}
}
