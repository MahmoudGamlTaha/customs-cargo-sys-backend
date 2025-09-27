package middleware

import (
	"fmt"
	"net/http"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"
	"Chumber-Workflow-System/pkg/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			c.Abort()
			return
		}

		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			c.Abort()
			return
		}
		fmt.Println("claims", claims)

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("claims", claims)
		c.Set("branch_id", claims.BranchID)
		c.Next()
	}
}

// RequireRole creates role-based authorization middleware
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		fmt.Println("role", role)
		if !exists {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
		c.Abort()
	}
}

// RequireAdmin creates admin-only authorization middleware
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

// RequireStaffOrAdmin creates staff or admin authorization middleware
func RequireStaffOrAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleStaff, models.RoleAdmin, models.RoleAuditor)
}

func RequireStaffOrAdminOrBranchAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleStaff, models.RoleAdmin, models.RoleAuditor,models.RoleBranchAdmin)
}

func RequireAccountantOrAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAccountant, models.RoleAdmin)
}

func RequireAccountantOrAdminOrBranchAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAccountant, models.RoleAdmin, models.RoleBranchAdmin)
}

// RequireClient creates client-only authorization middleware
func RequireClient() gin.HandlerFunc {
	return RequireRole(models.RoleClient)
}

// OptionalAuth creates optional authentication middleware
func OptionalAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("claims", claims)
		c.Set("branch_id", claims.BranchID)

		c.Next()
	}
}

// GetCurrentUserID gets the current user ID from context
func GetCurrentUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		panic("user_id not found in context")
	}
	return userID.(int64), true
}

// GetCurrentUserRole gets the current user role from context
func GetCurrentUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	return role.(string), true
}
func GetCurrentUserBranchID(c *gin.Context) (int64, bool) {
	branchID, exists := c.Get("branch_id")
	if !exists {
		return 0, false
	}
	return branchID.(int64), true
}

// GetCurrentUsername gets the current username from context
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == models.RoleAdmin
}

// IsStaff checks if the current user is staff
func IsStaff(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == models.RoleStaff
}

// IsClient checks if the current user is a client
func IsClient(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == models.RoleClient
}

// CanManageUsers checks if the current user can manage other users
func CanManageUsers(c *gin.Context) bool {
	return IsAdmin(c)
}

// CanManageRequests checks if the current user can manage requests
func CanManageRequests(c *gin.Context) bool {
	return IsStaff(c) || IsAdmin(c)
}

// RequirePasswordChangeMiddleware checks if user needs to change password
func RequirePasswordChangeMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip check for change-password route
		if c.Request.URL.Path == "/api/v1/users/change-password" {
			c.Next()
			return
		}

		// Get current user ID
		userID, exists := GetCurrentUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			c.Abort()
			return
		}

		// Get user details to check password reset flag
		user, err := userService.GetUser(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
			c.Abort()
			return
		}

		// Check if password reset is required
		if user.IsPasswordResetRequired {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "You must change your password first",
				Data:    nil,
				Error: &models.APIError{
					Code:    "PASSWORD_CHANGE_REQUIRED",
					Message: "You must change your password before accessing other resources",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
