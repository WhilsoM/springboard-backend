package repository

import (
	"context"
	"fmt"
	"log"
	"springboard/internal/opportunity/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OpportunityRepository interface {
	Create(ctx context.Context, employerID string, companyName string, data dto.CreateOpportunityRequest) (dto.OpportunityResponse, error)
	GetByID(ctx context.Context, id string) (dto.OpportunityResponse, error)
	Update(ctx context.Context, id string, data dto.CreateOpportunityRequest) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, f dto.SearchFilters) ([]dto.OpportunityResponse, error)
	GetByEmployerID(ctx context.Context, employerID string) ([]dto.OpportunityResponse, error)
}

type opportunityRepository struct {
	dbpool *pgxpool.Pool
}

func NewOpportunityRepository(dbpool *pgxpool.Pool) OpportunityRepository {
	return &opportunityRepository{dbpool: dbpool}
}

func (r *opportunityRepository) Create(ctx context.Context, employerID string, companyName string, data dto.CreateOpportunityRequest) (dto.OpportunityResponse, error) {
	var opp dto.OpportunityResponse
	query := `
		INSERT INTO opportunities (
			employer_id, company_name, title, description, type, format, city, address,
			latitude, longitude, tags, salary_min, salary_max, experience_level, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, employer_id, company_name, title, description, type, format, city, address,
		          latitude, longitude, tags, salary_min, salary_max, experience_level, is_active, expires_at, created_at`

	err := r.dbpool.QueryRow(ctx, query,
		employerID, companyName, data.Title, data.Description, data.Type, data.Format,
		data.City, data.Address, data.Latitude, data.Longitude, data.Tags,
		data.SalaryMin, data.SalaryMax, data.ExperienceLevel, data.ExpiresAt,
	).Scan(
		&opp.ID, &opp.EmployerID, &opp.CompanyName, &opp.Title, &opp.Description, &opp.Type, &opp.Format,
		&opp.City, &opp.Address, &opp.Latitude, &opp.Longitude, &opp.Tags,
		&opp.SalaryMin, &opp.SalaryMax, &opp.ExperienceLevel, &opp.IsActive, &opp.ExpiresAt, &opp.CreatedAt,
	)
	return opp, err
}

func (r *opportunityRepository) GetByID(ctx context.Context, id string) (dto.OpportunityResponse, error) {
	var o dto.OpportunityResponse
	query := `SELECT id, employer_id, company_name, title, description, type, format, city, address,
	                 latitude, longitude, tags, salary_min, salary_max, experience_level, is_active, expires_at, created_at
			  FROM opportunities WHERE id = $1`
	err := r.dbpool.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.EmployerID, &o.CompanyName, &o.Title, &o.Description, &o.Type, &o.Format,
		&o.City, &o.Address, &o.Latitude, &o.Longitude, &o.Tags,
		&o.SalaryMin, &o.SalaryMax, &o.ExperienceLevel, &o.IsActive, &o.ExpiresAt, &o.CreatedAt,
	)
	return o, err
}

func (r *opportunityRepository) Update(ctx context.Context, id string, data dto.CreateOpportunityRequest) error {
	query := `UPDATE opportunities
	          SET title=$1, description=$2, type=$3, format=$4, city=$5, address=$6,
			      latitude=$7, longitude=$8, tags=$9, salary_min=$10, salary_max=$11, experience_level=$12, expires_at=$13
			  WHERE id = $14`
	_, err := r.dbpool.Exec(ctx, query,
		data.Title, data.Description, data.Type, data.Format, data.City, data.Address,
		data.Latitude, data.Longitude, data.Tags, data.SalaryMin, data.SalaryMax, data.ExperienceLevel, data.ExpiresAt, id)
	return err
}

func (r *opportunityRepository) Delete(ctx context.Context, id string) error {
	_, err := r.dbpool.Exec(ctx, `DELETE FROM opportunities WHERE id = $1`, id)
	return err
}

func (r *opportunityRepository) GetByEmployerID(ctx context.Context, employerID string) ([]dto.OpportunityResponse, error) {
	query := `SELECT id, employer_id, company_name, title, description, type, format, city, address,
	                 latitude, longitude, tags, salary_min, salary_max, experience_level, is_active, expires_at, created_at
			  FROM opportunities WHERE employer_id = $1 ORDER BY created_at DESC`
	rows, err := r.dbpool.Query(ctx, query, employerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]dto.OpportunityResponse, 0)
	for rows.Next() {
		var o dto.OpportunityResponse
		rows.Scan(&o.ID, &o.EmployerID, &o.CompanyName, &o.Title, &o.Description, &o.Type, &o.Format,
			&o.City, &o.Address, &o.Latitude, &o.Longitude, &o.Tags,
			&o.SalaryMin, &o.SalaryMax, &o.ExperienceLevel, &o.IsActive, &o.ExpiresAt, &o.CreatedAt)
		res = append(res, o)
	}
	return res, nil
}

func (r *opportunityRepository) Search(ctx context.Context, f dto.SearchFilters) ([]dto.OpportunityResponse, error) {
	query := `SELECT id, employer_id, company_name, title, description, type, format, city, address,
	                 latitude, longitude, tags, salary_min, salary_max, experience_level, is_active, expires_at, created_at
			  FROM opportunities WHERE is_active = TRUE`
	args := []any{}
	argIdx := 1

	if f.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, f.Type)
		argIdx++
	}
	if f.Format != "" {
		query += fmt.Sprintf(" AND format = $%d", argIdx)
		args = append(args, f.Format)
		argIdx++
	}
	if f.City != "" {
		query += fmt.Sprintf(" AND city ILIKE $%d", argIdx)
		args = append(args, "%"+f.City+"%")
		argIdx++
	}
	if len(f.Tags) > 0 {
		query += fmt.Sprintf(" AND tags && $%d", argIdx)
		args = append(args, f.Tags)
		argIdx++
	}
	if f.Search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR company_name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+f.Search+"%")
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.dbpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]dto.OpportunityResponse, 0)
	for rows.Next() {
		var o dto.OpportunityResponse
		rows.Scan(&o.ID, &o.EmployerID, &o.CompanyName, &o.Title, &o.Description, &o.Type, &o.Format,
			&o.City, &o.Address, &o.Latitude, &o.Longitude, &o.Tags,
			&o.SalaryMin, &o.SalaryMax, &o.ExperienceLevel, &o.IsActive, &o.ExpiresAt, &o.CreatedAt)
		res = append(res, o)
	}
	log.Println("res ", res)
	return res, nil
}
