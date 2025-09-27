package service

import (
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
	"fmt"
)

// UserActivityService handles user activity business logic
type UserActivityService struct {
	repo     *repository.UserActivityRepository
	userRepo *repository.UserRepository
}

// NewUserActivityService creates a new user activity service
func NewUserActivityService(repo *repository.UserActivityRepository, userRepo *repository.UserRepository) *UserActivityService {
	return &UserActivityService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// GetAll retrieves all user activities with pagination based on user role
func (s *UserActivityService) GetAll(userID int64, pagination models.PaginationParams) ([]*models.UserActivities, int, error) {
	// Get the user to check their role
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting user details: %w", err)
	}

	// Check user role and filter accordingly
	if user.Role != nil && *user.Role == models.RoleAdmin {
		// Admin can see all activities
		return s.repo.GetAll(pagination)
	} else if user.Role != nil && *user.Role == models.RoleBranchAdmin {
		// Branch admin can only see activities from their branch
		if user.BranchID == nil {
			return nil, 0, fmt.Errorf("branch admin user has no branch assigned")
		}
		return s.repo.GetAllByBranch(*user.BranchID, pagination)
	} else {
		// Other users can only see their own activities
		return  s.repo.GetByUserID(userID, pagination)
	}
}

// GetByUserID retrieves user activities filtered by user ID
func (s *UserActivityService) GetByUserID(userID int64, pagination models.PaginationParams) ([]*models.UserActivities, int, error) {
	return s.repo.GetByUserID(userID, pagination)
}

// GetByModuleAndEntityID retrieves user activities filtered by module and entity ID
func (s *UserActivityService) GetByModuleAndEntityID(module string, entityID int64, pagination models.PaginationParams) ([]*models.UserActivities, int, error) {
	return s.repo.GetByModuleAndEntityID(module, entityID, pagination)
}
