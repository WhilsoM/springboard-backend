package repository

import (
	"context"
	"log"
	"springboard/internal/lib"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (lib.User, error)
	GetUserByID(ctx context.Context, id string) (lib.User, error)
	GetFullUserByID(ctx context.Context, id string, role lib.UserRole) (interface{}, error)
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
        SELECT id, email, role, display_name
        FROM users
        WHERE email = $1
    `

	err := r.dbpool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.DisplayName,
	)
	if err != nil {
		return lib.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (lib.User, error) {
	var user lib.User

	query := `
        SELECT id, email, role, display_name
        FROM users
        WHERE id = $1
    `

	err := r.dbpool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.DisplayName,
	)
	if err != nil {
		return lib.User{}, err
	}

	return user, nil
}

// get user by id with full data
func (r *userRepository) GetFullUserByID(ctx context.Context, id string, role lib.UserRole) (interface{}, error) {
	log.Print("GetFullUserByID repo start")

	switch role {
	case lib.RoleStudent:
		var applicant lib.ApplicantUser
		query := `
            SELECT u.id, u.email, u.role, u.display_name,
                   cp.university, cp.skills, cp.portfolio_url, cp.github_url
            FROM users u
            JOIN candidate_profiles cp ON u.id = cp.user_id
            WHERE u.id = $1`

		err := r.dbpool.QueryRow(ctx, query, id).Scan(
			&applicant.ID, &applicant.Email, &applicant.Role, &applicant.DisplayName,
			&applicant.University, &applicant.Skills, &applicant.PortfolioURL, &applicant.GithubURL,
		)
		if err != nil {
			log.Printf("Error scanning user profile role student: %v", err)
			return nil, err
		}
		return applicant, err

	case lib.RoleEmployer:
		var employer lib.EmployerUser
		query := `
            SELECT u.id, u.email, u.role, u.display_name,
                   ep.company_name, ep.is_verified, ep.inn
            FROM users u
            JOIN employer_profiles ep ON u.id = ep.user_id
            WHERE u.id = $1`

		err := r.dbpool.QueryRow(ctx, query, id).Scan(
			&employer.ID, &employer.Email, &employer.Role, &employer.DisplayName,
			&employer.CompanyName, &employer.IsVerified, &employer.INN,
		)
		log.Printf("Error scanning user profile role student: %v", err)

		return employer, err

	default:
		return r.GetUserByID(ctx, id)
	}
}
