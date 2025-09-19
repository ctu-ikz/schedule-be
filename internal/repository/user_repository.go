package repository

import (
	"context"
	"github.com/ctu-ikz/schedule-be/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
        INSERT INTO users (id, username, password)
        VALUES ($1, $2, $3)
        RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query, u.ID, u.Username, u.Password).
		Scan(&u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE username=$1`
	row := r.pool.QueryRow(ctx, query, username)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
