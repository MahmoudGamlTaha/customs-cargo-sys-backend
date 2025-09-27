package service

import (
	"database/sql"
	"errors"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
)

// BranchService handles branch business logic
type BranchService struct {
	branchRepo *repository.BranchRepository
}

// NewBranchService creates a new branch service
func NewBranchService(branchRepo *repository.BranchRepository) *BranchService {
	return &BranchService{
		branchRepo: branchRepo,
	}
}

// CreateBranch creates a new branch (admin only)
func (s *BranchService) CreateBranch(req *models.CreateBranchRequest) (*models.Branch, error) {
	// Check if branch code already exists
	exists, err := s.branchRepo.ExistsByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("branch code already exists")
	}

	// Check if branch name already exists
	exists, err = s.branchRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("branch name already exists")
	}

	// Create branch
	branch := &models.Branch{
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
	}

	err = s.branchRepo.Create(branch)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

// GetBranch retrieves a branch by ID
func (s *BranchService) GetBranch(id int64) (*models.Branch, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}

	return branch, nil
}

// GetBranchByCode retrieves a branch by code
func (s *BranchService) GetBranchByCode(code string) (*models.Branch, error) {
	branch, err := s.branchRepo.GetByCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}

	return branch, nil
}

// UpdateBranch updates a branch (admin only)
func (s *BranchService) UpdateBranch(id int64, req *models.UpdateBranchRequest) (*models.Branch, error) {
	// Get existing branch
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if new name already exists (excluding current branch)
		exists, err := s.branchRepo.ExistsByName(req.Name)
		if err != nil {
			return nil, err
		}
		if exists && req.Name != branch.Name {
			return nil, errors.New("branch name already exists")
		}
		branch.Name = req.Name
	}

	if req.Code != "" {
		// Check if new code already exists (excluding current branch)
		exists, err := s.branchRepo.ExistsByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if exists && req.Code != branch.Code {
			return nil, errors.New("branch code already exists")
		}
		branch.Code = req.Code
	}

	if req.Address != "" {
		branch.Address = req.Address
	}

	if req.Phone != "" {
		branch.Phone = req.Phone
	}

	if req.Email != "" {
		branch.Email = req.Email
	}

	// Update branch
	err = s.branchRepo.Update(branch)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

// ListBranches retrieves branches with pagination and filtering
func (s *BranchService) ListBranches(params models.FilterParams, pagination models.PaginationParams) ([]*models.Branch, *models.PaginationResponse, error) {
	branches, total, err := s.branchRepo.List(params, pagination)
	if err != nil {
		return nil, nil, err
	}

	paginationResp := models.CalculatePagination(pagination.Page, pagination.PageSize, total)
	return branches, &paginationResp, nil
}

// GetAllBranches retrieves all branches (for dropdown lists)
func (s *BranchService) GetAllBranches() ([]*models.Branch, error) {
	return s.branchRepo.GetAll()
}

// DeleteBranch deletes a branch (admin only)
func (s *BranchService) DeleteBranch(id int64) error {
	// Check if branch exists
	_, err := s.branchRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("branch not found")
		}
		return err
	}

	// TODO: Check if branch has associated users and handle accordingly
	// For now, we'll allow deletion (database constraints will prevent if there are references)

	return s.branchRepo.Delete(id)
}
