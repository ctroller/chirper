package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Find(username string) (*User, error)
}

type UserRepositoryImpl struct {
	DB *pgxpool.Pool
}

func (r *UserRepositoryImpl) Find(name string) (*User, error) {
	var userId int64
	var username, passwordHash string

	err := r.DB.QueryRow(context.Background(), "SELECT id, username, password_hash FROM users WHERE username=$1", name).Scan(&userId, &username, &passwordHash)

	if err != nil {
		return nil, err
	} else if userId == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &User{
		ID:           userId,
		Username:     username,
		PasswordHash: passwordHash,
	}, nil
}
