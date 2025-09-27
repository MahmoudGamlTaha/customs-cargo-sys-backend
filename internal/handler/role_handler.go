package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleBasedAccessHandler struct {
	service *service.RoleService
}

func NewRoleBasedAccessHandler(service *service.RoleService) *RoleBasedAccessHandler {
	return &RoleBasedAccessHandler{service: service}
}

// ListPermissions handles listing all permissions
// @Summary List permissions
// @Description List all permissions
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} models.APIResponse{data=models.Permission}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /admin/permissions [get]
func (h *RoleBasedAccessHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.service.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Permissions retrieved successfully", permissions))
}

// AssignPermissionsToRole assigns permissions to a role
// @Summary Assign permissions to a role
// @Description Assigns a list of permissions to a specific role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param permissions body []int64 true "List of Permission IDs"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /admin/roles/{id}/permissions [post]
func (h *RoleBasedAccessHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid role ID"))
		return
	}

	var permissionIDs []int64
	if err := c.ShouldBindJSON(&permissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	if err := h.service.AssignPermissionsToRole(roleID, permissionIDs); err != nil {
		if err.Error() == "role not found" || err.Error() == "one or more permissions not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse(err.Error()))
		} else if err.Error() == "super_admin role cannot be modified" {
			c.JSON(http.StatusForbidden, models.ConflictErrorResponse(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Permissions assigned successfully", nil))
}

// CreateRole handles the creation of a new role
// @Summary Create a new role
// @Description Create a new role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role body models.CreateRoleRequest true "Role object"
// @Success 201 {object} models.APIResponse{data=models.Role}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /admin/roles [post]
func (h *RoleBasedAccessHandler) CreateRole(c *gin.Context) {
	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(nil))
		return
	}

	createdRole, err := h.service.CreateRole(&req)
	if err != nil {
		if err.Error() == "role with this code already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
		} else if err.Error() == "role with this Arabic name already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
		} else if err.Error() == "role with this English name already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
		} else {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		}
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Role created", createdRole))
}

// ListRoles handles listing all roles
// @Summary List roles
// @Description List all roles
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.Role}
// @Failure 500 {object} models.APIResponse
// @Router /admin/roles [get]
func (h *RoleBasedAccessHandler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("Roles retrieved successfully", roles))
}

// ListUserPermissions handles listing all permissions for a user
// @Summary List user permissions
// @Description List all permissions for a specific user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.APIResponse{data=[]models.Permission}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /admin/users/{id}/permissions [get]
func (h *RoleBasedAccessHandler) ListUserPermissions(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid user ID"))
		return
	}

	permissions, err := h.service.ListUserPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User permissions retrieved successfully", permissions))
}

// GetCurrentUserPermissions handles listing all permissions for the current user
// @Summary List current user permissions
// @Description List all permissions for the currently authenticated user
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.Permission}
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /me/permissions [get]
func (h *RoleBasedAccessHandler) GetCurrentUserPermissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	permissions, err := h.service.ListUserPermissions(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("User permissions retrieved successfully", permissions))
}

// DeleteRole handles the deletion of a role
// @Summary Delete a role
// @Description Delete a role and detach it from users and permission_role tables
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /admin/roles/{id} [delete]
func (h *RoleBasedAccessHandler) DeleteRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid role ID"))
		return
	}

	if err := h.service.DeleteRole(roleID); err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse(err.Error()))
		} else if err.Error() == "super_admin role cannot be deleted" {
			c.JSON(http.StatusForbidden, models.ConflictErrorResponse(err.Error()))
		} else if err.Error() == "cannot delete role: users are currently assigned to this role" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Role deleted successfully", nil))
}

// GetRolePermissions handles listing all permissions for a specific role
// @Summary List role permissions
// @Description List all permissions for a specific role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} models.APIResponse{data=[]models.Permission}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /admin/roles/{id}/permissions [get]
func (h *RoleBasedAccessHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid role ID"))
		return
	}

	permissions, err := h.service.GetRolePermissions(roleID)
	if err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Role permissions retrieved successfully", permissions))
}
