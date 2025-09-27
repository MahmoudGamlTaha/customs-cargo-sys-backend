package service

import (
	"database/sql"
	"errors"
	"fmt"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
	"Chumber-Workflow-System/pkg/auth"
)

// UserService handles user business logic
type UserService struct {
	userRepo        *repository.UserRepository
	companyRepo     *repository.CompanyRepository
	branchRepo      *repository.BranchRepository
	roleRepo        *repository.RoleRepository
	passwordManager *auth.PasswordManager
	jwtManager      *auth.JWTManager
}

// NewUserService creates a new user service
func NewUserService(
	userRepo *repository.UserRepository,
	companyRepo *repository.CompanyRepository,
	branchRepo *repository.BranchRepository,
	roleRepo *repository.RoleRepository,
	passwordManager *auth.PasswordManager,
	jwtManager *auth.JWTManager,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		companyRepo:     companyRepo,
		branchRepo:      branchRepo,
		roleRepo:        roleRepo,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
	}
}

// Register registers a new user (client only)
func (s *UserService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Validate password strength
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Only allow client registration through public endpoint
	if req.Role == nil || *req.Role != models.RoleClient {
		return nil, errors.New("only client registration is allowed")
	}

	// Validate company if provided
	var companyID *int64
	tmpID := int64(req.CompanyID)
	companyID = &tmpID

	// Check if company exists
	_, err = s.companyRepo.GetByID(*companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	// Hash password
	passwordHash, err := s.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		FirstName:     req.FirstName,
		LastName:      &req.LastName,
		Phone:         req.Phone,
		Role:          req.Role,
		CompanyID:     companyID,
		IsActive:      true,
		EmailVerified: false, // In a real app, you'd send verification email
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

// Login authenticates a user and returns a token
func (s *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid username or password")
		}
		fmt.Println(err)
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user inactive")
	}

	// Verify password
	valid, err := s.passwordManager.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid username or password")
	}

	// Update last login
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login for user %s: %v\n", user.Username, err)
	}

	// If the user has a branch ID, fetch the branch details
	if user.BranchID != nil && *user.BranchID > 0 {
		branch, err := s.branchRepo.GetByID(*user.BranchID)
		if err == nil {
			user.Branch = branch
		} else {
			fmt.Printf("Failed to fetch branch details for user %s: %v\n", user.Username, err)
		}
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Email, *user.Role, *user.BranchID)
	if err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""

	return &models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// CreateUser creates a new user (admin only)
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Validate password strength
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Validate role
	if req.Role == nil || !models.IsValidRole(*req.Role) {
		return nil, errors.New("invalid role")
	}

	// Validate company for clients
	var companyID *int64
	if req.Role != nil && *req.Role == models.RoleClient && req.CompanyID != nil {
		// Parse string to int64
		companyIDVal := *req.CompanyID

		_, err = s.companyRepo.GetByID(companyIDVal)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("company not found")
			}
			return nil, err
		}
		companyID = &companyIDVal
	}

	// Validate branch for staff
	var branchID *int64
	if req.BranchID != nil {
		// Parse string to int64
		branchIDVal := *req.BranchID

		_, err = s.branchRepo.GetByID(branchIDVal)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("branch not found")
			}
			return nil, err
		}
		branchID = &branchIDVal
	}

	// Hash password
	passwordHash, err := s.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:                req.Username,
		Email:                   req.Email,
		PasswordHash:            passwordHash,
		FirstName:               req.FirstName,
		LastName:                &req.LastName,
		Phone:                   &req.Phone,
		Role:                    req.Role,
		RoleID:                  req.RoleID,
		CompanyID:               companyID,
		BranchID:                branchID,
		IsActive:                true,
		EmailVerified:           true, // Admin-created users are pre-verified
		IsPasswordResetRequired: true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int64) (*models.User, error) {
	// Repository expects int64 type
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		fmt.Println("user not found", err)
		return nil, err
	}

	// Clear password hash
	user.PasswordHash = ""
	return user, nil
}

// GetUserWithDetails retrieves a user with company and branch details
func (s *UserService) GetUserWithDetails(id int64) (*models.UserDetails, error) {
	// Repository expects int64 type
	userDetails, err := s.userRepo.GetWithDetails(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return userDetails, nil
}

// UpdateUser updates a user's profile
func (s *UserService) UpdateUser(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	// Repository expects int64 type
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = &req.LastName
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Email != "" {
		// Check if new email already exists
		if req.Email != user.Email {
			exists, err := s.userRepo.ExistsByEmail(req.Email)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("email already exists")
			}
		}
		user.Email = req.Email
	}

	// Update user
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user, _ = s.userRepo.GetByID(id)
	return user, nil
}

// UpdateUserRole updates a user's role and assignments (admin only)
func (s *UserService) UpdateUserRole(id int64, req *models.UpdateUserRoleRequest) (*models.User, error) {
	// Get existing user
	// Repository expects int64 type
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Validate role
	if req.Role == nil || !models.IsValidRole(*req.Role) {
		return nil, errors.New("invalid role")
	}

	// Update role
	user.Role = req.Role

	// Handle company assignment for clients
	if req.Role != nil && *req.Role == models.RoleClient && req.CompanyID != nil {
		var companyIDVal int64
		companyIDVal = *req.CompanyID

		_, err = s.companyRepo.GetByID(companyIDVal)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("company not found")
			}
			return nil, err
		}
		user.CompanyID = &companyIDVal
	} else {
		user.CompanyID = nil
	}

	// Handle branch assignment for staff
	if req.Role != nil && *req.Role == models.RoleStaff && req.BranchID != nil {
		var branchIDVal int64
		branchIDVal = *req.BranchID

		_, err = s.branchRepo.GetByID(branchIDVal)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("branch not found")
			}
			return nil, err
		}

		_, err = s.branchRepo.GetByID(branchIDVal)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("branch not found")
			}
			return nil, err
		}
		user.BranchID = &branchIDVal
	} else {
		user.BranchID = nil
	}

	// Update active status if provided
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update user
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id int64, req *models.ChangePasswordRequest) error {
	// Get user
	// Repository expects int64 type
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Verify current password
	valid, err := s.passwordManager.VerifyPassword(req.CurrentPassword, user.PasswordHash)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if err := auth.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	// Hash new password
	newPasswordHash, err := s.passwordManager.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password and set reset required flag to false
	err = s.userRepo.UpdatePasswordWithResetFlag(id, newPasswordHash, false)
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword resets a user's password (admin only)
func (s *UserService) ResetPassword(userID int64, req *models.ResetPasswordRequest) error {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Validate new password strength
	if err := auth.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	// Hash new password
	newPasswordHash, err := s.passwordManager.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password and set reset required flag to true
	err = s.userRepo.UpdatePasswordWithResetFlag(userID, newPasswordHash, true)
	if err != nil {
		return err
	}

	return nil
}

// ListUsers retrieves users with pagination and filtering
func (s *UserService) ListUsers(params models.FilterParams, pagination models.PaginationParams) ([]*models.User, *models.PaginationResponse, error) {
	// Repository expects int64 type
	users, total, err := s.userRepo.List(params, pagination)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	paginationResp := models.CalculatePagination(pagination.Page, pagination.PageSize, total)
	return users, &paginationResp, nil
}

// DeleteUser deletes a user (admin only)
func (s *UserService) DeleteUser(id int64) error {
	// Check if user exists
	// Repository expects int64 type
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Repository expects int64 type
	return s.userRepo.Delete(id)
}

// UpdateUserAdmin updates a user with admin privileges
func (s *UserService) UpdateUserAdmin(userID int64, req *models.UpdateUserAdminRequest) (*models.User, error) {
	// Get the existing user
	existingUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Prepare updates map with non-nil fields from the request
	updates := make(map[string]interface{})

	// Only update fields that are provided in the request
	if req.Username != nil {
		// Check if username already exists (if changed)
		if *req.Username != existingUser.Username {
			exists, err := s.userRepo.ExistsByUsername(*req.Username)
			if err != nil {
				return nil, fmt.Errorf("failed to check username availability: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("username already exists")
			}
		}
		updates["username"] = *req.Username
	}

	if req.Email != nil {
		// Check if email already exists (if changed)
		if *req.Email != existingUser.Email {
			exists, err := s.userRepo.ExistsByEmail(*req.Email)
			if err != nil {
				return nil, fmt.Errorf("failed to check email availability: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("email already exists")
			}
		}
		updates["email"] = *req.Email
	}

	if req.FirstName != nil {
		if *req.FirstName != existingUser.FirstName {
			updates["first_name"] = *req.FirstName
		}
	}

	if req.LastName != nil {
		if *req.LastName != *existingUser.LastName {
			updates["last_name"] = *req.LastName
		}
	}

	if req.Phone != nil {
		if *req.Phone != *existingUser.Phone {
			updates["phone"] = *req.Phone
		}
	}

	if req.Role != nil {
		// check if role is one of the valid roles
		if *req.Role != *existingUser.Role {
			if !models.IsValidRole(*req.Role) {
				return nil, fmt.Errorf("invalid role")
			}
			updates["role"] = *req.Role
		}
	}

	if req.CompanyID != nil {
		if *req.CompanyID != *existingUser.CompanyID {
			// Verify company exists if provided
			if *req.CompanyID != 0 {
				_, err := s.companyRepo.GetByID(*req.CompanyID)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return nil, fmt.Errorf("company not found")
					}
					return nil, fmt.Errorf("failed to verify company: %w", err)
				}
			}
			updates["company_id"] = req.CompanyID
		}
	}

	// If branch ID is being set, verify it exists
	if req.BranchID != nil && *req.BranchID != 0 {
		if existingUser.BranchID == nil || *req.BranchID != *existingUser.BranchID {
			_, err := s.branchRepo.GetByID(*req.BranchID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, fmt.Errorf("branch not found")
				}
				return nil, fmt.Errorf("failed to verify branch: %w", err)
			}
			updates["branch_id"] = req.BranchID
		}
	}

	if req.IsActive != nil {
		if *req.IsActive != existingUser.IsActive {
			updates["is_active"] = *req.IsActive
		}
	}

	// If no updates, return the existing user
	if len(updates) == 0 {
		return existingUser, nil
	}

	// Perform the update
	updatedUser, err := s.userRepo.PartialUpdate(userID, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// GetMemberCount returns the count of members based on user role and permissions
func (s *UserService) GetMemberCount(currentUser *models.User) (int64, error) {
	if currentUser == nil {
		return 0, errors.New("user not authenticated")
	}

	// Check if user is admin - can see all members
	if currentUser.IsAdmin() {
		return s.userRepo.CountAllMembers()
	}

	// Check if user is branch admin - can see members in their branch
	if currentUser.IsBranchAdmin() && currentUser.BranchID != nil {
		// For branch-level users, return count for their branch
		return s.userRepo.CountMembersByBranch(*currentUser.BranchID)
	}

	// Regular users cannot access member counts
	return 0, errors.New("insufficient permissions to view member count")
}

// GetCurrentUserRoleAndPermissions retrieves the current user's role and permissions
func (s *UserService) GetCurrentUserRoleAndPermissions(userID int64) (*models.CurrentUserRoleResponse, error) {
	// Get user details
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if user has a role assigned
	if user.RoleID == nil {
		return &models.CurrentUserRoleResponse{
			UserID:      user.ID,
			Username:    user.Username,
			Role:        nil,
			Permissions: []models.Permission{},
		}, nil
	}

	// Get role details
	role, err := s.roleRepo.GetByID(*user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// Get permissions for the role
	permissions, err := s.roleRepo.GetRolePermissions(*user.RoleID)
	if err != nil {
		return nil, err
	}

	return &models.CurrentUserRoleResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Role:        role,
		Permissions: permissions,
	}, nil
}
