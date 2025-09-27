package repository

import (
	"fmt"
	"strings"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"
)

// CompanyRepository handles company data operations
type CompanyRepository struct {
	db *database.DB
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *database.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Create creates a new company
func (r *CompanyRepository) Create(company *models.Company) error {
	query := `
		INSERT INTO companies (name, code, address, phone, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		company.Name,
		company.Code,
		company.Address,
		company.Phone,
		company.Email,
	).Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)

	return err
}

// GetByID retrieves a company by ID
func (r *CompanyRepository) GetByID(id int64) (*models.Company, error) {
	company := &models.Company{}
	query := `
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM companies 
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&company.ID,
		&company.Name,
		&company.Code,
		&company.Address,
		&company.Phone,
		&company.Email,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}

// GetByCode retrieves a company by code
func (r *CompanyRepository) GetByCode(code string) (*models.Company, error) {
	company := &models.Company{}
	query := `
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM companies 
		WHERE code = $1`

	err := r.db.QueryRow(query, code).Scan(
		&company.ID,
		&company.Name,
		&company.Code,
		&company.Address,
		&company.Phone,
		&company.Email,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}

// Update updates a company
func (r *CompanyRepository) Update(company *models.Company) error {
	query := `
		UPDATE companies 
		SET name = $2, code = $3, address = $4, phone = $5, email = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		company.ID,
		company.Name,
		company.Code,
		company.Address,
		company.Phone,
		company.Email,
	).Scan(&company.UpdatedAt)

	return err
}

// Delete deletes a company
func (r *CompanyRepository) Delete(id int64) error {
	query := `DELETE FROM companies WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves companies with pagination and filtering
func (r *CompanyRepository) List(params models.FilterParams, pagination models.PaginationParams) ([]*models.Company, int, error) {
	var companies []*models.Company
	var total int

	// Build WHERE clause
	whereClause, args := r.buildWhereClause(params)

	// Count total records
	countQuery := "SELECT COUNT(*) FROM companies" + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderClause := r.buildOrderClause(params)

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM companies%s%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, len(args)+1, len(args)+2)

	args = append(args, pagination.GetLimit(), pagination.GetOffset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		company := &models.Company{}
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Code,
			&company.Address,
			&company.Phone,
			&company.Email,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		companies = append(companies, company)
	}

	return companies, total, nil
}

// GetAll retrieves all companies (for dropdown lists)
func (r *CompanyRepository) GetAll() ([]*models.Company, error) {
	var companies []*models.Company

	query := `
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM companies
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		company := &models.Company{}
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Code,
			&company.Address,
			&company.Phone,
			&company.Email,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// ExistsByCode checks if a company exists by code
func (r *CompanyRepository) ExistsByCode(code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM companies WHERE code = $1)`
	err := r.db.QueryRow(query, code).Scan(&exists)
	return exists, err
}

// ExistsByName checks if a company exists by name
func (r *CompanyRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM companies WHERE name = $1)`
	err := r.db.QueryRow(query, name).Scan(&exists)
	return exists, err
}

// buildWhereClause builds the WHERE clause for filtering
func (r *CompanyRepository) buildWhereClause(params models.FilterParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR code ILIKE $%d OR address ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}

// buildOrderClause builds the ORDER BY clause
func (r *CompanyRepository) buildOrderClause(params models.FilterParams) string {
	sortBy := "name"
	if params.SortBy != "" {
		switch params.SortBy {
		case "name", "code", "created_at":
			sortBy = params.SortBy
		}
	}

	sortOrder := "ASC"
	if params.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}
