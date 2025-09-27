package middleware

import (
	"net/http"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"

	"github.com/gin-gonic/gin"
)

// HasPermission checks if the user has the required permission
func HasPermission(permissionRepo *repository.PermissionRepository, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetCurrentUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			c.Abort()
			return
		}

		hasPerm, err := permissionRepo.UserHasPermission(userID, permissionCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
			c.Abort()
			return
		}

		if !hasPerm {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			c.Abort()
			return
		}

		c.Next()
	}
}
