package service

import (
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
	"database/sql"
	"errors"
)

type RoleService struct {
	permissionRepo *repository.PermissionRepository
	roleRepo       *repository.RoleRepository
}

func NewRoleService(permissionRepo *repository.PermissionRepository, roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
	}
}

func (s *RoleService) CreateRole(req *models.CreateRoleRequest) (*models.Role, error) {
	// Check if a role with the same code already exists
	_, err := s.roleRepo.GetByCode(req.Code)
	if err == nil {
		return nil, errors.New("role with this code already exists")
	} else if err != sql.ErrNoRows {
		return nil, err // Other database error
	}

	// Check if a role with the same Arabic name already exists
	_, err = s.roleRepo.GetByNameAR(req.NameAR)
	if err == nil {
		return nil, errors.New("role with this Arabic name already exists")
	} else if err != sql.ErrNoRows {
		return nil, err // Other database error
	}

	// Check if a role with the same English name already exists
	_, err = s.roleRepo.GetByNameEN(req.NameEN)
	if err == nil {
		return nil, errors.New("role with this English name already exists")
	} else if err != sql.ErrNoRows {
		return nil, err // Other database error
	}

	role := &models.Role{
		NameAR: req.NameAR,
		NameEN: req.NameEN,
		Code:   req.Code,
	}
	return s.roleRepo.Create(role)
}

// AssignPermissionsToRole assigns permissions to a role
func (s *RoleService) AssignPermissionsToRole(roleID int64, permissionIDs []int64) error {
	// Validate that the role exists
	role, err := s.roleRepo.RoleExists(roleID)
	if err != nil {
		return err
	}
	if role.ID == 0 {
		return errors.New("role not found")
	}

	if role.Code == "super_admin" {
		return errors.New("super_admin role cannot be modified")
	}

	// Validate that all permissions exist
	permissionsExist, err := s.permissionRepo.PermissionsExist(permissionIDs)
	if err != nil {
		return err
	}
	if !permissionsExist {
		return errors.New("one or more permissions not found")
	}

	return s.roleRepo.AssignPermissionsToRole(roleID, permissionIDs)
}

// ListPermissions retrieves all permissions
func (s *RoleService) ListPermissions() ([]models.Permission, error) {
	return s.permissionRepo.ListPermissions()
}

// ListRoles retrieves all roles
func (s *RoleService) ListRoles() ([]models.Role, error) {
	return s.roleRepo.ListRoles()
}

// ListUserPermissions retrieves all permissions for a specific user
func (s *RoleService) ListUserPermissions(userID int64) ([]models.Permission, error) {
	return s.permissionRepo.ListUserPermissions(userID)
}

// DeleteRole deletes a role and detaches it from users and permission_role tables
func (s *RoleService) DeleteRole(roleID int64) error {
	// Validate that the role exists
	role, err := s.roleRepo.RoleExists(roleID)
	if err != nil {
		return err
	}
	if role.ID == 0 {
		return errors.New("role not found")
	}

	// Prevent deletion of super_admin role
	if role.Code == "super_admin" {
		return errors.New("super_admin role cannot be deleted")
	}

	// Check if any users have this role
	hasUsers, err := s.roleRepo.HasUsersWithRole(roleID)
	if err != nil {
		return err
	}
	if hasUsers {
		return errors.New("cannot delete role: users are currently assigned to this role")
	}

	return s.roleRepo.DeleteRole(roleID)
}

// GetRolePermissions retrieves all permissions for a specific role
func (s *RoleService) GetRolePermissions(roleID int64) ([]models.Permission, error) {
	// Validate that the role exists
	role, err := s.roleRepo.RoleExists(roleID)
	if err != nil {
		return nil, err
	}
	if role.ID == 0 {
		return nil, errors.New("role not found")
	}

	return s.roleRepo.GetRolePermissions(roleID)
}
