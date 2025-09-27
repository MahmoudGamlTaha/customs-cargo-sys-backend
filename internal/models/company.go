package models

import (
	
)

// Company represents a company entity
type Company struct {
	BaseModel
	Name    string `json:"name" db:"name" validate:"required,min=2,max=255"`
	Code    string `json:"code" db:"code" validate:"required,min=2,max=50"`
	Address string `json:"address" db:"address"`
	Phone   string `json:"phone" db:"phone" validate:"omitempty,min=10,max=20"`
	Email   string `json:"email" db:"email" validate:"omitempty,email"`
}

// CreateCompanyRequest represents the request payload for creating a company
type CreateCompanyRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Code    string `json:"code" validate:"required,min=2,max=50"`
	Address string `json:"address"`
	Phone   string `json:"phone" validate:"omitempty,min=10,max=20"`
	Email   string `json:"email" validate:"omitempty,email"`
}

// UpdateCompanyRequest represents the request payload for updating a company
type UpdateCompanyRequest struct {
	Name    string `json:"name" validate:"omitempty,min=2,max=255"`
	Code    string `json:"code" validate:"omitempty,min=2,max=50"`
	Address string `json:"address"`
	Phone   string `json:"phone" validate:"omitempty,min=10,max=20"`
	Email   string `json:"email" validate:"omitempty,email"`
}

// CompanyResponse represents the response payload for company data
type CompanyResponse struct {
	ID      int64 `json:"id"`
	Name    string    `json:"name"`
	Code    string    `json:"code"`
	Address string    `json:"address"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
}

// ToResponse converts Company to CompanyResponse
func (c *Company) ToResponse() *CompanyResponse {
	return &CompanyResponse{
		ID:      c.ID,
		Name:    c.Name,
		Code:    c.Code,
		Address: c.Address,
		Phone:   c.Phone,
		Email:   c.Email,
	}
}

