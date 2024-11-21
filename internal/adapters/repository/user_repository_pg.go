package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
)

type UserRepositoryPg struct {
	db *pgxpool.Pool
}

func NewUserRepositoryPg(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepositoryPg{db: db}
}

func (r *UserRepositoryPg) CreateUser(ctx context.Context, user *entities.User) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (r *UserRepositoryPg) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, email, password, created_at FROM users WHERE email=$1", email)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
