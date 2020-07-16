package user

type reader interface {
	FindById(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
}

type writer interface {
	Update(user *User) error
	Store(user *User) (string, error)
}

type Repository interface {
	reader
	writer
}
