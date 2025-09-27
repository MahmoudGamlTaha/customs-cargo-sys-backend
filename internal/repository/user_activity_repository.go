package repository

import (
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"
	"fmt"
)

// UserActivityRepositoryImpl implements the UserActivityRepository interface
type UserActivityRepository struct {
	db *database.DB
}

// NewUserActivityRepository creates a new instance of the user activity repository
func NewUserActivityRepository(db *database.DB) *UserActivityRepository {
	return &UserActivityRepository{db: db}
}

// Create inserts a new user activity record into the database
func (r *UserActivityRepository) Create(activity *models.UserActivities) error {
	query := `
		INSERT INTO user_activities (description, action, user_id, username, module, entity_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		activity.Description,
		activity.Action,
		activity.UserId,
		activity.Username,
		activity.Module,
		activity.EntityId,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

// GetAll retrieves all user activities with pagination
func (r *UserActivityRepository) GetAll(pagination models.PaginationParams) ([]*models.UserActivities, int, error) {
	// Count total records first
	var total int
	countQuery := `SELECT COUNT(*) FROM user_activities`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting user activities: %w", err)
	}

	// Then get the paginated records
	query := `
		SELECT id, description, action, user_id, username, module, entity_id, created_at
		FROM user_activities
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := r.db.Query(query, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying user activities: %w", err)
	}
	defer rows.Close()

	var activities []*models.UserActivities
	for rows.Next() {
		var activity models.UserActivities
		if err := rows.Scan(
			&activity.ID,
			&activity.Description,
			&activity.Action,
			&activity.UserId,
			&activity.Username,
			&activity.Module,
			&activity.EntityId,
			&activity.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning user activity: %w", err)
		}
		activities = append(activities, &activity)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return activities, total, nil
}

// GetAllByBranch retrieves all user activities filtered by branch ID with pagination
func (r *UserActivityRepository) GetAllByBranch(branchID int64, pagination models.PaginationParams) ([]*models.UserActivities, int, error) {
	// Count total records first
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM user_activities ua 
		JOIN users u ON ua.user_id = u.id 
		WHERE u.branch_id = $1
	`
	err := r.db.QueryRow(countQuery, branchID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting user activities by branch: %w", err)
	}

	// Then get the paginated records
	query := `
		SELECT ua.id, ua.description, ua.action, ua.user_id, ua.username, ua.module, ua.entity_id, ua.created_at
		FROM user_activities ua
		JOIN users u ON ua.user_id = u.id
		WHERE u.branch_id = $1
		ORDER BY ua.created_at DESC
		LIMIT $2 OFFSET $3
	`

	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := r.db.Query(query, branchID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying user activities by branch: %w", err)
	}
	defer rows.Close()

	var activities []*models.UserActivities
	for rows.Next() {
		var activity models.UserActivities
		if err := rows.Scan(
			&activity.ID,
			&activity.Description,
			&activity.Action,
			&activity.UserId,
			&activity.Username,
			&activity.Module,
			&activity.EntityId,
			&activity.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning user activity: %w", err)
		}
		activities = append(activities, &activity)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return activities, total, nil
}

// GetByUserID retrieves user activities filtered by user ID
func (r *UserActivityRepository) GetByUserID(userID int64, pagination models.PaginationParams) ([]*models.UserActivities, int, error) {
	// Count total records first
	var total int
	countQuery := `SELECT COUNT(*) FROM user_activities WHERE user_id = $1`
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting user activities: %w", err)
	}

	// Then get the paginated records
	query := `
		SELECT id, description, action, user_id, username, module, entity_id, created_at
		FROM user_activities
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := r.db.Query(query, userID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying user activities: %w", err)
	}
	defer rows.Close()

	var activities []*models.UserActivities
	for rows.Next() {
		var activity models.UserActivities
		if err := rows.Scan(
			&activity.ID,
			&activity.Description,
			&activity.Action,
			&activity.UserId,
			&activity.Username,
			&activity.Module,
			&activity.EntityId,
			&activity.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning user activity: %w", err)
		}
		activities = append(activities, &activity)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return activities, total, nil
}

// GetByModuleAndEntityID retrieves user activities filtered by module and entity ID
func (r *UserActivityRepository) GetByModuleAndEntityID(module string, entityID int64, pagination models.PaginationParams) ([]*models.UserActivities, int, error) {
	// Count total records first
	var total int
	countQuery := `SELECT COUNT(*) FROM user_activities WHERE module = $1 AND entity_id = $2`
	err := r.db.QueryRow(countQuery, module, entityID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting user activities: %w", err)
	}

	// Then get the paginated records
	query := `
		SELECT id, description, action, user_id, username, module, entity_id, created_at
		FROM user_activities
		WHERE module = $1 AND entity_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := r.db.Query(query, module, entityID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying user activities: %w", err)
	}
	defer rows.Close()

	var activities []*models.UserActivities
	for rows.Next() {
		var activity models.UserActivities
		if err := rows.Scan(
			&activity.ID,
			&activity.Description,
			&activity.Action,
			&activity.UserId,
			&activity.Username,
			&activity.Module,
			&activity.EntityId,
			&activity.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning user activity: %w", err)
		}
		activities = append(activities, &activity)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return activities, total, nil
}
