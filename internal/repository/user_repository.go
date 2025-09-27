package repository

import (
	"fmt"
	"strconv"
	"strings"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, phone, role, role_id, company_id, branch_id, is_active, email_verified, is_password_reset_required)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.RoleID,
		user.CompanyID,
		user.BranchID,
		user.IsActive,
		user.EmailVerified,
		user.IsPasswordResetRequired,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, phone, role, role_id, 
		       company_id, branch_id, is_active, email_verified, created_at, updated_at, last_login, is_password_reset_required

		FROM users 
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.RoleID,
		&user.CompanyID,
		&user.BranchID,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.IsPasswordResetRequired,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, phone, role, 
		       company_id, branch_id, is_active, email_verified, created_at, updated_at, last_login
		FROM users 
		WHERE username = $1`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.CompanyID,
		&user.BranchID,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, phone,role, role_id, 
		       company_id, branch_id, is_active, email_verified, created_at, updated_at, last_login, is_password_reset_required
		FROM users 
		WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.RoleID,
		&user.CompanyID,
		&user.BranchID,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.IsPasswordResetRequired,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetWithDetails retrieves a user with company and branch details
func (r *UserRepository) GetWithDetails(id int64) (*models.UserDetails, error) {
	userDetails := &models.UserDetails{}
	query := `
		SELECT id, username, email, first_name, last_name, phone, role, is_active, 
		       email_verified, created_at, last_login, company_name, company_code, 
		       branch_name, branch_code
		FROM user_details 
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&userDetails.ID,
		&userDetails.Username,
		&userDetails.Email,
		&userDetails.FirstName,
		&userDetails.LastName,
		&userDetails.Phone,
		&userDetails.Role,
		&userDetails.IsActive,
		&userDetails.EmailVerified,
		&userDetails.CreatedAt,
		&userDetails.LastLogin,
		&userDetails.CompanyName,
		&userDetails.CompanyCode,
		&userDetails.BranchName,
		&userDetails.BranchCode,
	)

	if err != nil {
		return nil, err
	}

	return userDetails, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, first_name = $4, last_name = $5, phone = $6, 
		    role = $7, company_id = $8, branch_id = $9, is_active = $10, email_verified = $11,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.CompanyID,
		user.BranchID,
		user.IsActive,
		user.EmailVerified,
	).Scan(&user.UpdatedAt)

	return err
}

// PartialUpdate updates specific fields of a user
func (r *UserRepository) PartialUpdate(userID int64, updates map[string]interface{}) (*models.User, error) {
	if len(updates) == 0 {
		return r.GetByID(userID)
	}

	setClauses := []string{}
	args := []interface{}{userID}
	argPos := 2 // $1 is for userID

	// Build the SET clause dynamically based on provided fields
	if username, ok := updates["username"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", argPos))
		args = append(args, username)
		argPos++
	}
	if email, ok := updates["email"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argPos))
		args = append(args, email)
		argPos++
	}
	if firstName, ok := updates["first_name"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("first_name = $%d", argPos))
		args = append(args, firstName)
		argPos++
	}
	if lastName, ok := updates["last_name"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("last_name = $%d", argPos))
		args = append(args, lastName)
		argPos++
	}
	if phone, ok := updates["phone"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, phone)
		argPos++
	}
	if role, ok := updates["role"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argPos))
		args = append(args, role)
		argPos++
	}

	if companyID, ok := updates["company_id"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("company_id = $%d", argPos))
		args = append(args, companyID)
		argPos++
	}
	if branchID, ok := updates["branch_id"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("branch_id = $%d", argPos))
		args = append(args, branchID)
		argPos++
	}
	if isActive, ok := updates["is_active"]; ok {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, isActive)
		argPos++
	}

	// Add updated_at to the set clauses
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	// Build the query
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $1 RETURNING id, username, email, first_name, last_name, phone, role_id, company_id, branch_id, is_active, created_at, updated_at, last_login, role",
		strings.Join(setClauses, ", "))

	// Execute the query
	user := &models.User{}
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.RoleID,
		&user.CompanyID,
		&user.BranchID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID int64, passwordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.Exec(query, userID, passwordHash)
	return err
}

// UpdatePasswordWithResetFlag updates password and sets is_password_reset_required flag
func (r *UserRepository) UpdatePasswordWithResetFlag(userID int64, passwordHash string, isResetRequired bool) error {
	query := `
		UPDATE users 
		SET password_hash = $2, is_password_reset_required = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.Exec(query, userID, passwordHash, isResetRequired)
	return err
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(userID int64) error {
	/*	query := `
			UPDATE users
			SET last_login = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1`
		fmt.Println(userID)
		_, err := r.db.Exec(query, userID)*/
	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves users with pagination and filtering
func (r *UserRepository) List(params models.FilterParams, pagination models.PaginationParams) ([]*models.User, int, error) {
	var users []*models.User
	var total int

	// Build WHERE clause
	whereClause, args := r.buildWhereClause(params)

	// Count total records
	countQuery := "SELECT COUNT(*) FROM users" + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderClause := r.buildOrderClause(params)

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, username, email, first_name, last_name, phone, role, is_active, 
		       email_verified, created_at, last_login, company_id, 
		       branch_id
		FROM users%s%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, len(args)+1, len(args)+2)

	args = append(args, pagination.GetLimit(), pagination.GetOffset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.LastLogin,
			&user.CompanyID,
			&user.BranchID,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// ExistsByUsername checks if a user exists by username
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

// ExistsByEmail checks if a user exists by email
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// buildWhereClause builds the WHERE clause for filtering
func (r *UserRepository) buildWhereClause(params models.FilterParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(username ILIKE $%d OR email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", argIndex, argIndex, argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	if params.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, params.Role)
		argIndex++
	}

	if params.CompanyID != "" {
		if companyID, err := strconv.ParseInt(params.CompanyID, 10, 64); err == nil {
			conditions = append(conditions, fmt.Sprintf("company_id = $%d", argIndex))
			args = append(args, companyID)
			argIndex++
		}
	}

	if params.BranchID != 0 {
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argIndex))
		args = append(args, params.BranchID)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}

// buildOrderClause builds the ORDER BY clause
func (r *UserRepository) buildOrderClause(params models.FilterParams) string {
	sortBy := "created_at"
	if params.SortBy != "" {
		switch params.SortBy {
		case "username", "email", "first_name", "last_name", "role", "created_at":
			sortBy = params.SortBy
		}
	}

	sortOrder := "DESC"
	if params.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}

// CountAllMembers counts all active users in the system
func (r *UserRepository) CountAllMembers() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// CountMembersByBranch counts active users in a specific branch
func (r *UserRepository) CountMembersByBranch(branchID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users WHERE branch_id = $1 AND is_active = true`
	err := r.db.QueryRow(query, branchID).Scan(&count)
	return count, err
}
