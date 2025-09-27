package service

import (
	"database/sql"
	"errors"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
)

// CompanyService handles company business logic
type CompanyService struct {
	companyRepo *repository.CompanyRepository
}

// NewCompanyService creates a new company service
func NewCompanyService(companyRepo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
	}
}

// CreateCompany creates a new company (admin only)
func (s *CompanyService) CreateCompany(req *models.CreateCompanyRequest) (*models.Company, error) {
	// Check if company code already exists
	exists, err := s.companyRepo.ExistsByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("company code already exists")
	}

	// Check if company name already exists
	exists, err = s.companyRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("company name already exists")
	}

	// Create company
	company := &models.Company{
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
	}

	err = s.companyRepo.Create(company)
	if err != nil {
		return nil, err
	}

	return company, nil
}

// GetCompany retrieves a company by ID
func (s *CompanyService) GetCompany(id int64) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	return company, nil
}

// GetCompanyByCode retrieves a company by code
func (s *CompanyService) GetCompanyByCode(code string) (*models.Company, error) {
	company, err := s.companyRepo.GetByCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	return company, nil
}

// UpdateCompany updates a company (admin only)
func (s *CompanyService) UpdateCompany(id int64, req *models.UpdateCompanyRequest) (*models.Company, error) {
	// Get existing company
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if new name already exists (excluding current company)
		exists, err := s.companyRepo.ExistsByName(req.Name)
		if err != nil {
			return nil, err
		}
		if exists && req.Name != company.Name {
			return nil, errors.New("company name already exists")
		}
		company.Name = req.Name
	}

	if req.Code != "" {
		// Check if new code already exists (excluding current company)
		exists, err := s.companyRepo.ExistsByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if exists && req.Code != company.Code {
			return nil, errors.New("company code already exists")
		}
		company.Code = req.Code
	}

	if req.Address != "" {
		company.Address = req.Address
	}

	if req.Phone != "" {
		company.Phone = req.Phone
	}

	if req.Email != "" {
		company.Email = req.Email
	}

	// Update company
	err = s.companyRepo.Update(company)
	if err != nil {
		return nil, err
	}

	return company, nil
}

// ListCompanies retrieves companies with pagination and filtering
func (s *CompanyService) ListCompanies(params models.FilterParams, pagination models.PaginationParams) ([]*models.Company, *models.PaginationResponse, error) {
	companies, total, err := s.companyRepo.List(params, pagination)
	if err != nil {
		return nil, nil, err
	}

	paginationResp := models.CalculatePagination(pagination.Page, pagination.PageSize, total)
	return companies, &paginationResp, nil
}

// GetAllCompanies retrieves all companies (for dropdown lists)
func (s *CompanyService) GetAllCompanies() ([]*models.Company, error) {
	return s.companyRepo.GetAll()
}

// DeleteCompany deletes a company (admin only)
func (s *CompanyService) DeleteCompany(id int64) error {
	// Check if company exists
	_, err := s.companyRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("company not found")
		}
		return err
	}

	// TODO: Check if company has associated users and handle accordingly
	// For now, we'll allow deletion (database constraints will prevent if there are references)

	return s.companyRepo.Delete(id)
}
