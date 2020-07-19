package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"

	"user-app/pkg/user"
	"user-app/pkg/user/repository"
)

const (
	driverName = "mysql"
)

type Repository struct {
	opts repository.MySQLOpts
	db   *sql.DB
}

// NewRepository create instance of Mysql repository
func NewRepository(opts repository.MySQLOpts) (*Repository, error) {
	db, err := sql.Open(driverName,
		fmt.Sprintf("%s:%s@%s(%s:%d)/%s", opts.User,
			opts.Password, opts.Protocol,
			opts.Host, opts.Port,
			opts.DatabaseName))
	if err != nil {
		return nil, err
	}
	for i := 1; i < 5; i++ {
		// Check connection with exponential back-off
		err = db.Ping()
		if err != nil {
			time.Sleep(time.Duration(1 << i) * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &Repository{
		opts: opts,
		db:   db,
	}, nil
}

// FindById finds user by id
func (r *Repository) FindById(ctx context.Context, id string) (*user.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	u := &user.User{}
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Email, &u.Name,
			&u.Address, &u.Telephone, &u.PasswordHash)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	tx.Commit()
	return nil, user.NotFound
}

// FindByEmail finds user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, "SELECT * FROM users WHERE Email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	u := &user.User{}
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Email, &u.Name,
			&u.Address, &u.Telephone, &u.PasswordHash)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	tx.Commit()
	return nil, user.NotFound
}

// Update update user entity
func (r *Repository) Update(ctx context.Context, u *user.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		"UPDATE users SET Name=?, Email=?, Telephone=?, Address=?, PasswordHash=? WHERE Id=?",
		u.Name, u.Email, u.Telephone, u.Address, u.PasswordHash, u.Id)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Store new user entity
func (r *Repository) Store(ctx context.Context, user *user.User) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO users(Id, Email, Name, Telephone, Address, PasswordHash) VALUES( ?, ?, ?, ?, ?, ? )",
		user.Id, user.Email, user.Name, user.Telephone, user.Address, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return user.Id, nil
}

// StoreToken store token entity
func (r *Repository) StoreToken(ctx context.Context, token *user.Token) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO tokens(Value, Email) VALUES( ?, ?)", token.Value, token.Email)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetByEmail get token by email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.Token, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, "SELECT * FROM tokens WHERE Email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	token := &user.Token{}
	for rows.Next() {
		err = rows.Scan(&token.Value, &token.Email)
		if err != nil {
			return nil, err
		}
		return token, nil
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return nil, user.NotFound
}

// GetByToken get token by it value
func (r *Repository) GetByToken(ctx context.Context, tokenValue string) (*user.Token, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, "SELECT Value, Email, Used FROM tokens WHERE Value = ?", tokenValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	token := &user.Token{}
	for rows.Next() {
		err = rows.Scan(&token.Value, &token.Email, &token.Used)
		if err != nil {
			return nil, err
		}
		return token, nil
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return nil, user.NotFound
}

// DisableToken mark token as used
func (r *Repository) DisableToken(ctx context.Context, tokenValue string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		"UPDATE tokens SET Used=? WHERE Value=?",
		false, tokenValue)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func init() {
	var _ repository.Repository = &Repository{}
}
