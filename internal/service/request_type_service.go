package service

import (
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
)

// RequestTypeService handles business logic for request types
type RequestTypeService struct {
	requestTypeRepo *repository.RequestTypeRepository
}

// NewRequestTypeService creates a new request type service
func NewRequestTypeService(requestTypeRepo *repository.RequestTypeRepository) *RequestTypeService {
	return &RequestTypeService{requestTypeRepo: requestTypeRepo}
}

// ListRequestTypes retrieves request types with pagination and filtering
func (s *RequestTypeService) ListRequestTypes(params models.FilterParams, pagination models.PaginationParams) ([]*models.RequestType, *models.PaginationResponse, error) {
	items, total, err := s.requestTypeRepo.List(params, pagination)
	if err != nil {
		return nil, nil, err
	}
	paginationResp := models.CalculatePagination(pagination.Page, pagination.PageSize, total)
	return items, &paginationResp, nil
}

// GetAllRequestTypes retrieves all request types (for dropdown lists)
func (s *RequestTypeService) GetAllRequestTypes() ([]*models.RequestType, error) {
	return s.requestTypeRepo.GetAll()
}
