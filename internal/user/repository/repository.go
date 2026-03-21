package repository

import (
	"context"
	"fmt"
	"log"
	"springboard/internal/lib"
	"springboard/internal/user/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (lib.User, error)
	GetUserByID(ctx context.Context, id string) (lib.User, error)
	GetFullUserByID(ctx context.Context, id string, role lib.UserRole) (any, error)
	DeleteUserByID(ctx context.Context, id string) error
	UpdateCandidate(ctx context.Context, userID string, data dto.UpdateMeCandidateRequest) error
	UpdateEmployer(ctx context.Context, userID string, data dto.UpdateMeEmployerRequest) error
	VerifyEmployer(ctx context.Context, userID string, inn string) error
	UpdatePrivacy(ctx context.Context, userID string, isPrivate bool) error
	UpdateAvatar(ctx context.Context, userID string, role lib.UserRole, url string) error
	GetApplicants(ctx context.Context, query string, limit, offset int) ([]lib.ApplicantUser, error)
	CreateContactRequest(ctx context.Context, senderID, receiverID string) error
	UpdateContactStatus(ctx context.Context, requestID, status string) error
	GetContacts(ctx context.Context, userID string) ([]lib.User, error)
}

type userRepository struct {
	dbpool *pgxpool.Pool
}

func NewUserRepository(dbpool *pgxpool.Pool) UserRepository {
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
func (r *userRepository) GetFullUserByID(ctx context.Context, id string, role lib.UserRole) (any, error) {
	log.Print("GetFullUserByID repo start")

	switch role {
	case lib.RoleStudent:
		var applicant lib.ApplicantUser
		query := `
						SELECT u.id, u.email, u.role, u.display_name,
                   COALESCE(cp.university, ''),
                   COALESCE(cp.course, 0),
                   COALESCE(cp.skills, '{}'),
                   COALESCE(cp.portfolio_url, ''),
                   COALESCE(cp.github_url, ''),
                   COALESCE(cp.updated_at, NOW())
            FROM users u
            LEFT JOIN candidate_profiles cp ON u.id = cp.user_id
            WHERE u.id = $1`

		err := r.dbpool.QueryRow(ctx, query, id).Scan(
			&applicant.ID, &applicant.Email, &applicant.Role, &applicant.DisplayName,
			&applicant.University, &applicant.Course, &applicant.Skills, &applicant.PortfolioURL, &applicant.GithubURL,
			&applicant.UpdatedAt,
		)

		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}
		return applicant, nil

	case lib.RoleEmployer:
		var employer lib.EmployerUser
		query := `
						SELECT u.id, u.email, u.role, u.display_name,
                   COALESCE(ep.company_name, ''),
                   COALESCE(ep.is_verified, false),
                   COALESCE(ep.inn, ''),
                   COALESCE(ep.description, ''),
                   COALESCE(ep.website_url, ''),
                   COALESCE(ep.logo_url, ''),
                   COALESCE(ep.updated_at, NOW())
            FROM users u
            LEFT JOIN employer_profiles ep ON u.id = ep.user_id
            WHERE u.id = $1`

		err := r.dbpool.QueryRow(ctx, query, id).Scan(
			&employer.ID, &employer.Email, &employer.Role, &employer.DisplayName,
			&employer.CompanyName, &employer.IsVerified, &employer.INN,
			&employer.Description, &employer.WebsiteURL, &employer.LogoURL,
			&employer.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}
		return employer, nil

	default:
		return r.GetUserByID(ctx, id)
	}
}

func (r *userRepository) DeleteUserByID(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.dbpool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

func (r *userRepository) UpdateCandidate(ctx context.Context, id string, data dto.UpdateMeCandidateRequest) error {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryUser := `UPDATE users SET display_name = $1, email = $2, is_private = $3 WHERE id = $4`
	_, err = tx.Exec(ctx, queryUser, data.DisplayName, data.Email, data.IsPrivate, id)
	if err != nil {
		return err
	}

	queryProfile := `
        UPDATE candidate_profiles
        SET university = $1, course = $2, skills = $3, portfolio_url = $4, github_url = $5, updated_at = NOW()
        WHERE user_id = $6`
	_, err = tx.Exec(ctx, queryProfile,
		data.University, data.Course, data.Skills, data.PortfolioURL, data.GithubURL, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *userRepository) UpdateEmployer(ctx context.Context, userID string, data dto.UpdateMeEmployerRequest) error {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryUser := `
        UPDATE users
        SET email = $1, display_name = $2, is_private = $3
        WHERE id = $4`

	_, err = tx.Exec(ctx, queryUser, data.Email, data.DisplayName, data.IsPrivate, userID)
	if err != nil {
		return fmt.Errorf("failed to update users table: %w", err)
	}

	queryProfile := `
        UPDATE employer_profiles
        SET company_name = $1, description = $2, website_url = $3, logo_url = $4, is_verified = $5, updated_at = NOW()
        WHERE user_id = $6`

	_, err = tx.Exec(ctx, queryProfile,
		data.CompanyName,
		data.Description,
		data.WebsiteURL,
		data.LogoURL,
		data.IsVerified,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update employer_profiles table: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *userRepository) VerifyEmployer(ctx context.Context, userID string, inn string) error {
	query := `UPDATE employer_profiles SET inn = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.dbpool.Exec(ctx, query, inn, userID)
	return err
}

func (r *userRepository) UpdatePrivacy(ctx context.Context, userID string, isPrivate bool) error {
	query := `UPDATE users SET is_private = $1 WHERE id = $2`
	_, err := r.dbpool.Exec(ctx, query, isPrivate, userID)
	return err
}

func (r *userRepository) UpdateAvatar(ctx context.Context, userID string, role lib.UserRole, url string) error {
	var query string
	if role == lib.RoleEmployer {
		query = `UPDATE employer_profiles SET logo_url = $1, updated_at = NOW() WHERE user_id = $2`
	} else {
		query = `UPDATE candidate_profiles SET avatar_url = $1, updated_at = NOW() WHERE user_id = $2`
	}
	_, err := r.dbpool.Exec(ctx, query, url, userID)
	return err
}
func (r *userRepository) GetApplicants(ctx context.Context, query string, limit, offset int) ([]lib.ApplicantUser, error) {
	sql := `
        SELECT u.id, u.email, u.role, u.display_name,
               COALESCE(cp.university, ''),
               COALESCE(cp.skills, '{}'),
               COALESCE(cp.avatar_url, '')
        FROM users u
        JOIN candidate_profiles cp ON u.id = cp.user_id
        WHERE u.role = 'applicant'
          AND u.is_private = false
          AND (u.display_name ILIKE $1 OR cp.university ILIKE $1)
        LIMIT $2 OFFSET $3`

	rows, err := r.dbpool.Query(ctx, sql, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query applicants: %w", err)
	}
	defer rows.Close()

	applicants := make([]lib.ApplicantUser, 0, limit)
	for rows.Next() {
		var a lib.ApplicantUser
		err := rows.Scan(&a.ID, &a.Email, &a.Role, &a.DisplayName, &a.University, &a.Skills, &a.AvatarURL)
		if err != nil {
			return nil, fmt.Errorf("scan applicant: %w", err)
		}
		applicants = append(applicants, a)
	}
	return applicants, nil
}

func (r *userRepository) CreateContactRequest(ctx context.Context, senderID, receiverID string) error {
	query := `INSERT INTO contacts (sender_id, receiver_id, status) VALUES ($1, $2, 'pending')`
	_, err := r.dbpool.Exec(ctx, query, senderID, receiverID)
	return err
}

func (r *userRepository) UpdateContactStatus(ctx context.Context, requestID, status string) error {
	query := `UPDATE contacts SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.dbpool.Exec(ctx, query, status, requestID)
	return err
}

func (r *userRepository) GetContacts(ctx context.Context, userID string) ([]lib.User, error) {
	query := `
        SELECT u.id, u.email, u.role, u.display_name
        FROM users u
        JOIN contacts c ON (c.sender_id = u.id OR c.receiver_id = u.id)
        WHERE (c.sender_id = $1 OR c.receiver_id = $1)
          AND c.status = 'accepted'
          AND u.id != $1`

	rows, err := r.dbpool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := make([]lib.User, 0)
	for rows.Next() {
		var u lib.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.DisplayName); err != nil {
			return nil, err
		}
		contacts = append(contacts, u)
	}
	return contacts, nil
}
