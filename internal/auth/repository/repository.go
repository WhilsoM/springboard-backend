package repository

import (
	"context"
	"springboard/internal/lib"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user lib.User) (lib.User, error)
	LoginUser(ctx context.Context, email string) (lib.User, error)
}

type authRepository struct {
	dbpool *pgxpool.Pool
}

func NewAuthRepository(dbpool *pgxpool.Pool) AuthRepository {
	return &authRepository{
		dbpool: dbpool,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user lib.User) (lib.User, error) {
	var createdUser lib.User

	query := `
        INSERT INTO users (email, password_hash, role, display_name)
        VALUES ($1, $2, $3, $4)
        RETURNING id, role, display_name, email
    `

	err := r.dbpool.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.DisplayName,
	).Scan(&createdUser.ID, &createdUser.Role, &createdUser.DisplayName, &createdUser.Email)
	if err != nil {
		return lib.User{}, err
	}

	return createdUser, nil
}

func (r *authRepository) LoginUser(ctx context.Context, email string) (lib.User, error) {
	var user lib.User
	query := `SELECT id, email, password_hash, role, display_name FROM users WHERE email = $1`

	err := r.dbpool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.DisplayName,
	)
	return user, err
}
