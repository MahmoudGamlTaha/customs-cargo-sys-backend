package repository

import (
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"
	"database/sql"
	"time"
)

type RoleRepository struct {
	db *database.DB
}

func NewRoleRepository(db *database.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.Role) (*models.Role, error) {
	query := `
		INSERT INTO roles (name_ar, name_en, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	err := r.db.QueryRow(query, role.NameAR, role.NameEN, role.Code, time.Now(), time.Now()).Scan(&role.ID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) AssignRoleToUser(userID int64, roleID int64) error {
	query := `
		UPDATE users 
		SET role_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2;
	`
	_, err := r.db.Exec(query, roleID, userID)
	return err
}

func (r *RoleRepository) GetByCode(code string) (*models.Role, error) {
	query := `SELECT id, name_ar, name_en, code FROM roles WHERE code = $1;`
	var role models.Role
	err := r.db.QueryRow(query, code).Scan(&role.ID, &role.NameAR, &role.NameEN, &role.Code)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetByID(id int64) (*models.Role, error) {
	query := `SELECT id, name_ar, name_en, code FROM roles WHERE id = $1;`
	var role models.Role
	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.NameAR, &role.NameEN, &role.Code)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) ListRoles() ([]models.Role, error) {
	query := `SELECT id, name_ar, name_en, code FROM roles;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.NameAR, &role.NameEN, &role.Code); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// RoleExists checks if a role exists in the database
func (r *RoleRepository) RoleExists(roleID int64) (models.Role, error) {
	query := `SELECT id, name_ar, name_en, code FROM roles WHERE id = $1;`
	var role models.Role
	err := r.db.QueryRow(query, roleID).Scan(&role.ID, &role.NameAR, &role.NameEN, &role.Code)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Role{}, nil
		}
		return models.Role{}, err
	}
	return role, nil
}

// AssignPermissionsToRole assigns a list of permissions to a role in a transaction
func (r *RoleRepository) AssignPermissionsToRole(roleID int64, permissionIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// First, remove all existing permissions for this role
	_, err = tx.Exec(`DELETE FROM permission_role WHERE role_id = $1;`, roleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then, add the new permissions
	stmt, err := tx.Prepare(`INSERT INTO permission_role (role_id, permission_id) VALUES ($1, $2);`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, permissionID := range permissionIDs {
		_, err = stmt.Exec(roleID, permissionID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// DeleteRole deletes a role and detaches it from users and permission_role tables
func (r *RoleRepository) DeleteRole(roleID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// First, set users with this role to have no role (NULL)
	_, err = tx.Exec(`UPDATE users SET role_id = NULL, updated_at = CURRENT_TIMESTAMP WHERE role_id = $1;`, roleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Second, remove all permissions associated with this role
	_, err = tx.Exec(`DELETE FROM permission_role WHERE role_id = $1;`, roleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Finally, delete the role itself
	_, err = tx.Exec(`DELETE FROM roles WHERE id = $1;`, roleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetRolePermissions retrieves all permissions for a specific role
func (r *RoleRepository) GetRolePermissions(roleID int64) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name_ar, p.name_en, p.code, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN permission_role pr ON p.id = pr.permission_id
		WHERE pr.role_id = $1
		ORDER BY p.name_en;
	`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		if err := rows.Scan(&permission.ID, &permission.NameAR, &permission.NameEN, &permission.Code, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

// HasUsersWithRole checks if any users have the specified role
func (r *RoleRepository) HasUsersWithRole(roleID int64) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE role_id = $1`
	err := r.db.QueryRow(query, roleID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByNameAR retrieves a role by Arabic name
func (r *RoleRepository) GetByNameAR(nameAR string) (*models.Role, error) {
	query := `SELECT id, name_ar, name_en, code FROM roles WHERE name_ar = $1;`
	var role models.Role
	err := r.db.QueryRow(query, nameAR).Scan(&role.ID, &role.NameAR, &role.NameEN, &role.Code)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByNameEN retrieves a role by English name
func (r *RoleRepository) GetByNameEN(nameEN string) (*models.Role, error) {
	query := `SELECT id, name_ar, name_en, code FROM roles WHERE name_en = $1;`
	var role models.Role
	err := r.db.QueryRow(query, nameEN).Scan(&role.ID, &role.NameAR, &role.NameEN, &role.Code)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
