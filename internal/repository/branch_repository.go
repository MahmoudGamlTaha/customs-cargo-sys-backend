package repository

import (
	"fmt"
	"strings"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"
)

// BranchRepository handles branch data operations
type BranchRepository struct {
	db *database.DB
}

// NewBranchRepository creates a new branch repository
func NewBranchRepository(db *database.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

// Create creates a new branch
func (r *BranchRepository) Create(branch *models.Branch) error {
	query := `
		INSERT INTO branches (name, code, address, phone, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		branch.Name,
		branch.Code,
		branch.Address,
		branch.Phone,
		branch.Email,
	).Scan(&branch.ID, &branch.CreatedAt, &branch.UpdatedAt)

	return err
}

// GetByID retrieves a branch by ID
func (r *BranchRepository) GetByID(id int64) (*models.Branch, error) {
	branch := &models.Branch{}
	query := `
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM branches 
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&branch.ID,
		&branch.Name,
		&branch.Code,
		&branch.Address,
		&branch.Phone,
		&branch.Email,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return branch, nil
}

// GetByCode retrieves a branch by code
func (r *BranchRepository) GetByCode(code string) (*models.Branch, error) {
	branch := &models.Branch{}
	query := `
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM branches 
		WHERE code = $1`

	err := r.db.QueryRow(query, code).Scan(
		&branch.ID,
		&branch.Name,
		&branch.Code,
		&branch.Address,
		&branch.Phone,
		&branch.Email,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return branch, nil
}

// Update updates a branch
func (r *BranchRepository) Update(branch *models.Branch) error {
	query := `
		UPDATE branches 
		SET name = $2, code = $3, address = $4, phone = $5, email = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		branch.ID,
		branch.Name,
		branch.Code,
		branch.Address,
		branch.Phone,
		branch.Email,
	).Scan(&branch.UpdatedAt)

	return err
}

// Delete deletes a branch
func (r *BranchRepository) Delete(id int64) error {
	query := `DELETE FROM branches WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves branches with pagination and filtering
func (r *BranchRepository) List(params models.FilterParams, pagination models.PaginationParams) ([]*models.Branch, int, error) {
	var branches []*models.Branch
	var total int

	// Build WHERE clause
	whereClause, args := r.buildWhereClause(params)

	// Count total records
	countQuery := "SELECT COUNT(*) FROM branches" + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderClause := r.buildOrderClause(params)

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM branches%s%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, len(args)+1, len(args)+2)

	args = append(args, pagination.GetLimit(), pagination.GetOffset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		branch := &models.Branch{}
		err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Code,
			&branch.Address,
			&branch.Phone,
			&branch.Email,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		branches = append(branches, branch)
	}

	return branches, total, nil
}

// GetAll retrieves all branches (for dropdown lists)
func (r *BranchRepository) GetAll() ([]*models.Branch, error) {
	var branches []*models.Branch

	query := `
		SELECT id, name, code, address, phone, email, created_at, updated_at
		FROM branches
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		branch := &models.Branch{}
		err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Code,
			&branch.Address,
			&branch.Phone,
			&branch.Email,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

// ExistsByCode checks if a branch exists by code
func (r *BranchRepository) ExistsByCode(code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM branches WHERE code = $1)`
	err := r.db.QueryRow(query, code).Scan(&exists)
	return exists, err
}

// ExistsByName checks if a branch exists by name
func (r *BranchRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM branches WHERE name = $1)`
	err := r.db.QueryRow(query, name).Scan(&exists)
	return exists, err
}

// buildWhereClause builds the WHERE clause for filtering
func (r *BranchRepository) buildWhereClause(params models.FilterParams) (string, []interface{}) {
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
func (r *BranchRepository) buildOrderClause(params models.FilterParams) string {
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
