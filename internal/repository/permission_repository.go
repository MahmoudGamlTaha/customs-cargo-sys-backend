package repository

import (
	"database/sql"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"

	"github.com/lib/pq"
)

// PermissionRepository handles database operations for permissions
type PermissionRepository struct {
	db *database.DB
}

// NewPermissionRepository creates a new PermissionRepository
func NewPermissionRepository(db *database.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// UserHasPermission checks if a user has a specific permission
// ListPermissions retrieves all permissions from the database
func (r *PermissionRepository) ListPermissions() ([]models.Permission, error) {
	query := `SELECT id, name_ar, name_en, code FROM permissions;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.NameAR, &p.NameEN, &p.Code); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// UserHasPermission checks if a user has a specific permission
// PermissionsExist checks if a list of permission IDs exist in the database
func (r *PermissionRepository) PermissionsExist(permissionIDs []int64) (bool, error) {
	query := `SELECT COUNT(id) FROM permissions WHERE id = ANY($1);`
	var count int
	err := r.db.QueryRow(query, pq.Array(permissionIDs)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == len(permissionIDs), nil
}

// ListUserPermissions retrieves all permissions for a specific user
func (r *PermissionRepository) ListUserPermissions(userID int64) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name_ar, p.name_en, p.code, p.created_at, p.updated_at
		FROM permissions p
		JOIN permission_role pr ON p.id = pr.permission_id
		JOIN users u ON pr.role_id = u.role_id
		WHERE u.id = $1;
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.NameAR, &p.NameEN, &p.Code, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

func (r *PermissionRepository) UserHasPermission(userID int64, permissionCode string) (bool, error) {
	query := `
		SELECT 1
		FROM users u
		JOIN permission_role pr ON u.role_id = pr.role_id
		JOIN permissions p ON pr.permission_id = p.id
		WHERE u.id = $1 AND p.code = $2
		LIMIT 1;
	`
	var exists int
	err := r.db.QueryRow(query, userID, permissionCode).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No permission found, but not an error
		}
		return false, err // An actual error occurred
	}
	return true, nil
}
