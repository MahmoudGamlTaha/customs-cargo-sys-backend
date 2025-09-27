package repository

import (
	"fmt"
	"strings"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"
)

// RequestTypeRepository handles request type data operations
type RequestTypeRepository struct {
	db *database.DB
}

// NewRequestTypeRepository creates a new request type repository
func NewRequestTypeRepository(db *database.DB) *RequestTypeRepository {
	return &RequestTypeRepository{db: db}
}

// List retrieves request types with pagination and filtering
func (r *RequestTypeRepository) List(params models.FilterParams, pagination models.PaginationParams) ([]*models.RequestType, int, error) {
	var items []*models.RequestType
	var total int

	whereClause, args := r.buildWhereClause(params)

	// Count total
	countQuery := "SELECT COUNT(*) FROM request_types" + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderClause := r.buildOrderClause(params)

	query := fmt.Sprintf(`
		SELECT id, name, name_en, price, ratio, created_at, updated_at
		FROM request_types%s%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, len(args)+1, len(args)+2)

	args = append(args, pagination.GetLimit(), pagination.GetOffset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		rt := &models.RequestType{}
		if err := rows.Scan(
			&rt.ID,
			&rt.Name,
			&rt.NameEn,
			&rt.Price,
			&rt.Ratio,
			&rt.CreatedAt,
			&rt.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, rt)
	}

	return items, total, nil
}

// GetAll retrieves all request types (for dropdown lists)
func (r *RequestTypeRepository) GetAll() ([]*models.RequestType, error) {
	var items []*models.RequestType

	query := `
		SELECT id, name, name_en, price, ratio, created_at, updated_at
		FROM request_types
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rt := &models.RequestType{}
		if err := rows.Scan(
			&rt.ID,
			&rt.Name,
			&rt.NameEn,
			&rt.Price,
			&rt.Ratio,
			&rt.CreatedAt,
			&rt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, rt)
	}

	return items, nil
}

// GetByID retrieves a request type by ID
func (r *RequestTypeRepository) GetByID(id int64) (*models.RequestType, error) {
	rt := &models.RequestType{}
	query := `
		SELECT id, name, name_en, price, ratio, created_at, updated_at
		FROM request_types
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&rt.ID,
		&rt.Name,
		&rt.NameEn,
		&rt.Price,
		&rt.Ratio,
		&rt.CreatedAt,
		&rt.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return rt, nil
}

func (r *RequestTypeRepository) buildWhereClause(params models.FilterParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR name_en ILIKE $%d)", argIndex, argIndex+1))
		like := "%" + params.Search + "%"
		args = append(args, like, like)
		argIndex += 2
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}
	return where, args
}

func (r *RequestTypeRepository) buildOrderClause(params models.FilterParams) string {
	sortBy := "name"
	if params.SortBy != "" {
		switch params.SortBy {
		case "name", "price", "created_at":
			sortBy = params.SortBy
		}
	}

	sortOrder := "ASC"
	if strings.ToUpper(params.SortOrder) == "DESC" {
		sortOrder = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}
