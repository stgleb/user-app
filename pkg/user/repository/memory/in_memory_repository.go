package memory

import (
	"context"
	"github.com/google/uuid"
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

func (r *InMemoryRepository) FindById(_ context.Context, id string) (*user.User, error) {
	if _, ok := r.userDb[id]; !ok {
		return nil, user.NotFound
	}
	return r.userDb[id], nil
}

func (r *InMemoryRepository) FindByEmail(_ context.Context, email string) (*user.User, error) {
	for _, u := range r.userDb {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, user.NotFound
}

func (r *InMemoryRepository) Update(_ context.Context, user *user.User) error {
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	r.userDb[user.Id] = user
	return nil
}

func (r *InMemoryRepository) Store(_ context.Context, user *user.User) (string, error) {
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	r.userDb[user.Id] = user
	return user.Id, nil
}

func (r *InMemoryRepository) StoreToken(_ context.Context, token *user.Token) error {
	r.tokenDb[token.Email] = token
	return nil
}

func (r *InMemoryRepository) GetByEmail(_ context.Context, email string) (*user.Token, error) {
	token, ok := r.tokenDb[email]
	if !ok {
		return nil, user.NotFound
	}
	return token, nil
}

func (r *InMemoryRepository) GetByToken(_ context.Context, tokenValue string) (*user.Token, error) {
	for _, token := range r.tokenDb {
		if token.Value == tokenValue {
			return token, nil
		}
	}
	return nil, user.NotFound
}

func (r *InMemoryRepository) DisableToken(_ context.Context, tokenValue string) error {
	for _, token := range r.tokenDb {
		if token.Value == tokenValue {
			token.Used = true
			return nil
		}
	}
	return user.NotFound
}

func init() {
	var _ repository.Repository = &InMemoryRepository{}
}
