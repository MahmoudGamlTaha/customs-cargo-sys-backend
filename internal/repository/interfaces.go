package repository

import (
	"Chumber-Workflow-System/internal/models"
)

// BaseRequestRepository defines the interface for request-related database operations
type BaseRequestRepository interface {
	Create(request *models.Request) error
	CreateRequestDetail(detail *models.RequestDetail, requestID int64) error
	GetByID(id int64) (*models.Request, error)
	GetWithDetails(id int64) (*models.RequestDetailDto, error)
	Update(request *models.Request) error
	Delete(id int64) error
	List(params models.FilterParams, pagination models.PaginationParams) ([]*models.RequestDetail, int, error)
	ListByClientID(clientID int64, params models.FilterParams, pagination models.PaginationParams) ([]*models.RequestDetail, int, error)
	GetHistory(requestID int64) ([]*models.RequestHistory, error)
	Approve(id int64, approverID int64) error
	Reject(id int64, approvedByStaffID int64, rejectionReason string) error
	MarkAsPaid(id int64, approvedByStaffID int64) error
}
