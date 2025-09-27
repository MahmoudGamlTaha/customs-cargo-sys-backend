package models

import (
	"time"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError represents an API error
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int `json:"page" form:"page" validate:"min=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// FilterParams represents common filter parameters
type FilterParams struct {
	Search         string `json:"search" form:"search"`
	Status         string `json:"status" form:"status"`
	Role           string `json:"role" form:"role"`
	CompanyID      string `json:"company_id" form:"company_id"`
	BranchID       int64  `json:"branch_id" form:"branch_id"`
	StartDate      string `json:"start_date" form:"start_date"`
	EndDate        string `json:"end_date" form:"end_date"`
	SortBy         string `json:"sort_by" form:"sort_by"`
	SortOrder      string `json:"sort_order" form:"sort_order"`
	Serial         string `json:"serial" form:"serial"`
	BranchIDFilter *int64 `json:"branch_id_filter" form:"branch_id_filter"`
}

// SuccessResponse creates a successful API response
func SuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// ErrorResponse creates an error API response
func ErrorResponse(code, message string, details map[string]interface{}) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(details map[string]interface{}) *APIResponse {
	return ErrorResponse("VALIDATION_ERROR", "Validation failed", details)
}

// UnauthorizedErrorResponse creates an unauthorized error response
func UnauthorizedErrorResponse() *APIResponse {
	return ErrorResponse("UNAUTHORIZED", "Authentication required", nil)
}

// ForbiddenErrorResponse creates a forbidden error response
func ForbiddenErrorResponse() *APIResponse {
	return ErrorResponse("FORBIDDEN", "Access denied", nil)
}

// NotFoundErrorResponse creates a not found error response
func NotFoundErrorResponse(resource string) *APIResponse {
	return ErrorResponse("NOT_FOUND", resource+" not found", nil)
}

// InternalErrorResponse creates an internal server error response
func InternalErrorResponse() *APIResponse {
	return ErrorResponse("INTERNAL_ERROR", "Internal server error", nil)
}

// ConflictErrorResponse creates a conflict error response
func ConflictErrorResponse(message string) *APIResponse {
	return ErrorResponse("CONFLICT", message, nil)
}

// BadRequestErrorResponse creates a bad request error response
func BadRequestErrorResponse(message string) *APIResponse {
	return ErrorResponse("BAD_REQUEST", message, nil)
}

// GetDefaultPagination returns default pagination parameters
func GetDefaultPagination() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 20,
	}
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize, total int) PaginationResponse {
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// GetOffset calculates the offset for database queries
func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the page size as limit
func (p PaginationParams) GetLimit() int {
	return p.PageSize
}
