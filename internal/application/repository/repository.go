package repository

import (
	"context"
	"errors"
	"springboard/internal/application/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadyApplied = errors.New("already applied")

type ApplicationRepository interface {
	Create(ctx context.Context, applicantID, oppID string, req dto.ApplyRequest) (dto.ApplicationResponse, error)
	GetByApplicant(ctx context.Context, applicantID string, limit, offset int) ([]dto.ApplicationResponse, error)
	GetByOpportunity(ctx context.Context, oppID string, limit, offset int) ([]dto.ApplicationResponse, error)
	UpdateStatus(ctx context.Context, appID, status string) error
	CheckOwnership(ctx context.Context, appID, employerID string) (bool, error)
}

type applicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, applicantID, oppID string, req dto.ApplyRequest) (dto.ApplicationResponse, error) {
	var res dto.ApplicationResponse
	err := r.db.QueryRow(ctx,
		`INSERT INTO applications (opportunity_id, applicant_id, cover_letter)
		 VALUES ($1, $2, $3) RETURNING id, opportunity_id, status, created_at`,
		oppID, applicantID, req.CoverLetter).Scan(&res.ID, &res.OpportunityID, &res.Status, &res.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return res, ErrAlreadyApplied
		}
		return res, err
	}
	return res, nil
}

func (r *applicationRepository) GetByOpportunity(ctx context.Context, oppID string, limit, offset int) ([]dto.ApplicationResponse, error) {
	query := `
		SELECT a.id, a.opportunity_id, a.status, a.cover_letter, a.created_at,
		       u.id, u.display_name, COALESCE(cp.university, ''), COALESCE(cp.course, 0),
		       COALESCE(cp.skills, '{}'), COALESCE(cp.github_url, ''), COALESCE(cp.avatar_url, '')
		FROM applications a
		JOIN users u ON a.applicant_id = u.id
		LEFT JOIN candidate_profiles cp ON u.id = cp.user_id
		WHERE a.opportunity_id = $1
		ORDER BY a.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, oppID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []dto.ApplicationResponse
	for rows.Next() {
		var a dto.ApplicationResponse
		var p dto.ApplicantInfo
		err := rows.Scan(&a.ID, &a.OpportunityID, &a.Status, &a.CoverLetter, &a.CreatedAt,
			&p.ID, &p.DisplayName, &p.University, &p.Course, &p.Skills, &p.GitHubURL, &p.AvatarURL)
		if err != nil {
			return nil, err
		}
		a.Applicant = &p
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *applicationRepository) GetByApplicant(ctx context.Context, applicantID string, limit, offset int) ([]dto.ApplicationResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, opportunity_id, status, created_at FROM applications
		 WHERE applicant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		applicantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dto.ApplicationResponse])
}

func (r *applicationRepository) UpdateStatus(ctx context.Context, appID, status string) error {
	_, err := r.db.Exec(ctx, "UPDATE applications SET status = $1, updated_at = NOW() WHERE id = $2", status, appID)
	return err
}

func (r *applicationRepository) CheckOwnership(ctx context.Context, appID, employerID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM applications a
		 JOIN opportunities o ON a.opportunity_id = o.id
		 WHERE a.id = $1 AND o.employer_id = $2)`, appID, employerID).Scan(&exists)
	return exists, err
}
