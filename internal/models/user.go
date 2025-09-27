package models

import (
	"time"
)

// User represents a user entity
type User struct {
	BaseModel
	Username                string     `json:"username" db:"username" validate:"required,min=3,max=100"`
	Email                   string     `json:"email" db:"email" validate:"required,email"`
	PasswordHash            string     `json:"-" db:"password_hash"`
	FirstName               string     `json:"first_name" db:"first_name" validate:"required,min=2,max=100"`
	LastName                *string    `json:"last_name" db:"last_name" validate:"required,min=2,max=100"`
	Phone                   *string    `json:"phone" db:"phone" validate:"omitempty,min=10,max=20"`
	RoleID                  *int64     `json:"role_id" db:"role_id" validate:"required"`
	CompanyID               *int64     `json:"company_id" db:"company_id"`
	BranchID                *int64     `json:"branch_id" db:"branch_id"`
	IsActive                bool       `json:"is_active" db:"is_active"`
	EmailVerified           bool       `json:"email_verified" db:"email_verified"`
	IsPasswordResetRequired bool       `json:"is_password_reset_required" db:"is_password_reset_required"`
	LastLogin               *time.Time `json:"last_login" db:"last_login"`
	Role                    *string    `json:"role" db:"role"`
	// Relationships (populated when needed)
	Company *Company `json:"company,omitempty"`
	Branch  *Branch  `json:"branch,omitempty"`
}

// UserDetails represents user with company and branch details
type UserDetails struct {
	ID                      int64      `json:"id" db:"id"`
	Username                string     `json:"username" db:"username"`
	Email                   string     `json:"email" db:"email"`
	FirstName               string     `json:"first_name" db:"first_name"`
	LastName                string     `json:"last_name" db:"last_name"`
	Phone                   string     `json:"phone" db:"phone"`
	Role                    *string    `json:"role" db:"role"`
	IsActive                bool       `json:"is_active" db:"is_active"`
	EmailVerified           bool       `json:"email_verified" db:"email_verified"`
	IsPasswordResetRequired bool       `json:"is_password_reset_required" db:"is_password_reset_required"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	LastLogin               *time.Time `json:"last_login" db:"last_login"`
	CompanyName             string     `json:"company_name" db:"company_name"`
	CompanyCode             string     `json:"company_code" db:"company_code"`
	BranchName              string     `json:"branch_name" db:"branch_name"`
	BranchCode              string     `json:"branch_code" db:"branch_code"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=100"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6,max=100"`
	FirstName string  `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=100"`
	Phone     *string `json:"phone" validate:"omitempty,min=10,max=20"`
	Role      *string `json:"role" validate:"required"`
	CompanyID int     `json:"company_id"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

// CreateUserRequest represents the request payload for creating a user (admin only)
type CreateUserRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=100"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6,max=100"`
	FirstName string  `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=100"`
	Phone     string  `json:"phone" validate:"omitempty,min=10,max=20"`
	Role      *string `json:"role" validate:"required"`
	CompanyID *int64  `json:"company_id" validate:"omitempty,gt=0"`
	BranchID  *int64  `json:"branch_id" validate:"omitempty,gt=0"`
	RoleID    *int64  `json:"role_id" validate:"omitempty,gt=0"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
	Email     string `json:"email" validate:"omitempty,email"`
}

// UpdateUserRoleRequest represents the request payload for updating user role (admin only)
type UpdateUserRoleRequest struct {
	Role      *string `json:"role" validate:"required"`
	CompanyID *int64  `json:"company_id" validate:"omitempty,gt=0"`
	BranchID  *int64  `json:"branch_id" validate:"omitempty,gt=0"`
	IsActive  *bool   `json:"is_active"`
}

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
}

// ResetPasswordRequest represents the request payload for admin password reset
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}

// UpdateUserAdminRequest represents the request payload for admin user updates
type UpdateUserAdminRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=100"`
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Phone     *string `json:"phone" validate:"omitempty,min=10,max=20"`
	CompanyID *int64  `json:"company_id" validate:"omitempty,gt=0"`
	BranchID  *int64  `json:"branch_id" validate:"omitempty,gt=0"`
	IsActive  *bool   `json:"is_active"`
	Role      *string `json:"role"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID                      int64      `json:"id"`
	Username                string     `json:"username"`
	Email                   string     `json:"email"`
	FirstName               string     `json:"first_name"`
	LastName                string     `json:"last_name"`
	Phone                   string     `json:"phone"`
	IsActive                bool       `json:"is_active"`
	EmailVerified           bool       `json:"email_verified"`
	IsPasswordResetRequired bool       `json:"is_password_reset_required"`
	CreatedAt               time.Time  `json:"created_at"`
	LastLogin               *time.Time `json:"last_login"`
	Company                 *Company   `json:"company,omitempty"`
	BranchID                *int64     `json:"branch_id,omitempty"`
	Branch                  *Branch    `json:"branch,omitempty"`
	RoleName                *string    `json:"role_name"`
	Role                    *string    `json:"role"`
	RoleID                  *int64     `json:"role_id"`
}

func (u *UserResponse) ToResponse() *UserResponse {
	panic("unimplemented")
}

// LoginResponse represents the response payload for login
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// CurrentUserRoleResponse represents the response payload for current user role and permissions
type CurrentUserRoleResponse struct {
	UserID      int64        `json:"user_id"`
	Username    string       `json:"username"`
	Role        *Role        `json:"role"`
	Permissions []Permission `json:"permissions"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	lastName := ""
	if u.LastName != nil {
		lastName = *u.LastName
	}
	phone := ""
	if u.Phone != nil {
		phone = *u.Phone
	}

	return &UserResponse{
		ID:                      u.ID,
		Username:                u.Username,
		Email:                   u.Email,
		FirstName:               u.FirstName,
		LastName:                lastName,
		Phone:                   phone,
		IsActive:                u.IsActive,
		EmailVerified:           u.EmailVerified,
		IsPasswordResetRequired: u.IsPasswordResetRequired,
		CreatedAt:               u.CreatedAt,
		LastLogin:               u.LastLogin,
		Company:                 u.Company,
		Branch:                  u.Branch,
		RoleName:                u.Role,
		Role:                    u.Role,
		RoleID:                  u.RoleID,
		BranchID:                u.BranchID,
	}
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	if u.LastName == nil || *u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + *u.LastName
}

// IsClient checks if the user is a client
func (u *User) IsClient() bool {
	return u.Role != nil && *u.Role == RoleClient
}

// IsStaff checks if the user is a staff member
func (u *User) IsStaff() bool {
	return u.Role != nil && *u.Role == RoleStaff
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role != nil && *u.Role == RoleAdmin
}

// IsBranchAdmin checks if the user is a branch admin
func (u *User) IsBranchAdmin() bool {
	return u.Role != nil && *u.Role == RoleBranchAdmin
}

// CanManageUsers checks if the user can manage other users
func (u *User) CanManageUsers() bool {
	return u.Role != nil && *u.Role == RoleAdmin
}

// CanManageRequests checks if the user can manage requests
func (u *User) CanManageRequests() bool {
	return u.Role != nil && (*u.Role == RoleStaff || *u.Role == RoleAdmin)
}
