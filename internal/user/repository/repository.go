package repository

import (
	"context"
	"springboard/internal/lib"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (lib.User, error)
	GetUserByID(ctx context.Context, id string) (lib.User, error)
}

type userRepository struct {
	dbpool *pgxpool.Pool
}

func NewUserRepository(dbpool *pgxpool.Pool) *userRepository {
	return &userRepository{
		dbpool: dbpool,
	}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (lib.User, error) {
	var user lib.User

	query := `
        SELECT id, email, password_hash, role, full_name
        FROM users
        WHERE email = $1
    `

	err := r.dbpool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.FullName,
	)
	if err != nil {
		return lib.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (lib.User, error) {
	var user lib.User

	query := `
        SELECT id, email, password_hash, role, full_name
        FROM users
        WHERE id = $1
    `

	err := r.dbpool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.FullName,
	)
	if err != nil {
		return lib.User{}, err
	}

	return user, nil
}
