package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"Chumber-Workflow-System/internal/middleware"
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new client user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("User registered successfully", user.ToResponse()))
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	loginResponse, err := h.userService.Login(&req)
	if err != nil {
		if err.Error() == "invalid username or password" || err.Error() == "account is deactivated" {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Login successful", loginResponse))
}

// GetProfile handles getting current user profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	userID := userIDStr

	user, err := h.userService.GetUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User profile retrieved", user.ToResponse()))
}

// UpdateProfile handles updating current user profile
// @Summary Update current user profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateUserRequest true "Update user request"
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	userID := userIDStr

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Profile updated successfully", user.ToResponse()))
}

// ChangePassword handles changing user password
// @Summary Change user password
// @Description Change the password of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	userID := userIDStr

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	err := h.userService.ChangePassword(userID, &req)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Password changed successfully", nil))
}

// CreateUser handles creating a new user (admin only)
// @Summary Create a new user
// @Description Create a new user (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateUserRequest true "Create user request"
// @Success 201 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /admin/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		fmt.Println(err)
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("User created successfully", user.ToResponse()))
}

// GetUser handles getting a user by ID (admin only)
// @Summary Get user by ID
// @Description Get a user by ID (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid user ID"))
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User retrieved", user.ToResponse()))
}

// ListUsers handles listing users with pagination (admin only)
// @Summary List users
// @Description List users with pagination and filtering (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Param role query string false "Filter by role"
// @Param company_id query string false "Filter by company ID"
// @Param branch_id query string false "Filter by branch ID"
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	pagination := models.GetDefaultPagination()
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}
	branchIDStr := c.Query("branch_id")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)
	if err != nil {
		branchID = 0
	}
	// Parse filter parameters
	params := models.FilterParams{
		Search:    c.Query("search"),
		Role:      c.Query("role"),
		CompanyID: c.Query("company_id"),
		BranchID:  branchID,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	users, paginationResp, err := h.userService.ListUsers(params, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	response := map[string]interface{}{
		"users":      users,
		"pagination": paginationResp,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Users retrieved", response))
}

// UpdateUserRole handles updating a user's role (admin only)
// @Summary Update user role
// @Description Update a user's role and assignments (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRoleRequest true "Update user role request"
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid user ID"))
		return
	}

	var req models.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	user, err := h.userService.UpdateUserRole(userID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User role updated successfully", user.ToResponse()))
}

// DeleteUser handles deleting a user (admin only)
// @Summary Delete user
// @Description Delete a user (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid user ID"))
		return
	}

	err = h.userService.DeleteUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User deleted successfully", nil))
}

// ResetPassword handles resetting a user's password (admin only)
// @Summary Reset user password
// @Description Reset a user's password (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body models.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/users/{id}/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid user ID"))
		return
	}

	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	err = h.userService.ResetPassword(userID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Password reset successfully", nil))
}


// UpdateUserAdmin handles updating a user with admin privileges
// @Summary Update user (admin only)
// @Description Update a user with admin privileges (can update any field except password)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body models.UpdateUserAdminRequest true "Update user request"
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUserAdmin(c *gin.Context) {
	// Get user ID from URL
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid user ID"))
		return
	}

	// Bind request body
	var req models.UpdateUserAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	// Call service to update user
	user, err := h.userService.UpdateUserAdmin(userID, &req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
		case "username already exists":
			c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Username already exists"))
		case "email already exists":
			c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Email already exists"))
		case "company not found":
			c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Company not found"))
		case "branch not found":
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Branch"))
		case "branch does not belong to the specified company":
			c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Branch does not belong to the specified company"))
		case "invalid role":
			c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid role"))
		default:
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User updated successfully", user.ToResponse()))
}

// GetMemberCount handles getting member count based on user permissions
// @Summary Get member count
// @Description Get count of members - all members for admin, branch members for branch admin
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=map[string]int64}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/member-count [get]
func (h *UserHandler) GetMemberCount(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	// Get member count based on user permissions
	count, err := h.userService.GetMemberCount(user)
	if err != nil {
		switch err.Error() {
		case "user not authenticated":
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		case "insufficient permissions to view member count":
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
		default:
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		}
		return
	}

	response := map[string]int64{
		"member_count": count,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Member count retrieved successfully", response))
}

// GetCurrentUserRoleAndPermissions handles getting current user's role and permissions
// @Summary Get current user role and permissions
// @Description Get the role and permissions of the currently authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.CurrentUserRoleResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/role-permissions [get]
func (h *UserHandler) GetCurrentUserRoleAndPermissions(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	rolePermissions, err := h.userService.GetCurrentUserRoleAndPermissions(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("User"))
			return
		}
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Role"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User role and permissions retrieved successfully", rolePermissions))
}
