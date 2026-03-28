package repository

import (
	"context"
	"fmt"
	"log"
	"springboard/internal/admin/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository interface {
	CreateCurator(ctx context.Context, email, hash, displayName string) error
	CreateTag(ctx context.Context, name string) (dto.Tag, error)
	GetAllTags(ctx context.Context) ([]dto.Tag, error)
	CreateVerification(ctx context.Context, employerID, inn, companyName, status string) error
	UpdateVerificationStatus(ctx context.Context, requestID, status string) error
	ApproveEmployer(ctx context.Context, employerID string) error
	DeleteOpportunity(ctx context.Context, oppID string) error
}

type adminRepository struct {
	dbpool *pgxpool.Pool
}

func NewAdminRepository(dbpool *pgxpool.Pool) AdminRepository {
	return &adminRepository{dbpool: dbpool}
}

func (r *adminRepository) CreateCurator(ctx context.Context, email, hash, displayName string) error {
	query := `INSERT INTO users (email, password_hash, display_name, role) VALUES ($1, $2, $3, 'curator')`
	_, err := r.dbpool.Exec(ctx, query, email, hash, displayName)
	return err
}

func (r *adminRepository) CreateTag(ctx context.Context, name string) (dto.Tag, error) {
	var tag dto.Tag
	query := `INSERT INTO tags (name) VALUES ($1) RETURNING id, name`
	err := r.dbpool.QueryRow(ctx, query, name).Scan(&tag.ID, &tag.Name)
	return tag, err
}

func (r *adminRepository) GetAllTags(ctx context.Context) ([]dto.Tag, error) {
	query := `SELECT id, name FROM tags ORDER BY name ASC`
	rows, err := r.dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []dto.Tag
	for rows.Next() {
		var t dto.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *adminRepository) CreateVerification(ctx context.Context, employerID, inn, companyName, status string) error {
	query := `INSERT INTO verifications (employer_id, inn, company_name, status) VALUES ($1, $2, $3, $4)`
	_, err := r.dbpool.Exec(ctx, query, employerID, inn, companyName, status)
	return err
}

func (r *adminRepository) UpdateVerificationStatus(ctx context.Context, requestID, status string) error {
	query := `UPDATE verifications SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.dbpool.Exec(ctx, query, status, requestID)
	return err
}

func (r *adminRepository) ApproveEmployer(ctx context.Context, employerID string) error {
	query := `UPDATE employer_profiles SET is_verified = true, updated_at = NOW() WHERE user_id = $1`
	_, err := r.dbpool.Exec(ctx, query, employerID)
	return err
}

func (r *adminRepository) DeleteOpportunity(ctx context.Context, oppID string) error {
	query := `DELETE FROM opportunities WHERE id = $1`
	result, err := r.dbpool.Exec(ctx, query, oppID)
	if err != nil {
		return err
	}
	log.Println("result", result)
	if result.RowsAffected() == 0 {
		return fmt.Errorf("opportunity not found")
	}
	return nil
}
