package repository

import (
	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u *model.User) error {
	query := `
    INSERT INTO users (email, password_hash, name, role)
    VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.Exec(query, u.Email, u.PasswordHash, u.Name, u.Role)
	return err
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var u model.User
	query := `
    SELECT id, email, password_hash, name, role, created_at, updated_at
    FROM users
    WHERE email=$1
    `
	if err := r.db.Get(&u, query, email); err != nil {
		return nil, err
	}
	return &u, nil
}
