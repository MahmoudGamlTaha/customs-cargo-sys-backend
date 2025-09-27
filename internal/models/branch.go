package models

import (
	
)

// Branch represents a branch entity
type Branch struct {
	BaseModel
	Name    string `json:"name" db:"name" validate:"required,min=2,max=255"`
	Code    string `json:"code" db:"code" validate:"required,min=2,max=50"`
	Address string `json:"address" db:"address"`
	Phone   string `json:"phone" db:"phone" validate:"omitempty,min=10,max=20"`
	Email   string `json:"email" db:"email" validate:"omitempty,email"`
}

// CreateBranchRequest represents the request payload for creating a branch
type CreateBranchRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Code    string `json:"code" validate:"required,min=2,max=50"`
	Address string `json:"address"`
	Phone   string `json:"phone" validate:"omitempty,min=10,max=20"`
	Email   string `json:"email" validate:"omitempty,email"`
}

// UpdateBranchRequest represents the request payload for updating a branch
type UpdateBranchRequest struct {
	Name    string `json:"name" validate:"omitempty,min=2,max=255"`
	Code    string `json:"code" validate:"omitempty,min=2,max=50"`
	Address string `json:"address"`
	Phone   string `json:"phone" validate:"omitempty,min=10,max=20"`
	Email   string `json:"email" validate:"omitempty,email"`
}

// BranchResponse represents the response payload for branch data
type BranchResponse struct {
	ID      int64 `json:"id"`
	Name    string    `json:"name"`
	Code    string    `json:"code"`
	Address string    `json:"address"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
}

// ToResponse converts Branch to BranchResponse
func (b *Branch) ToResponse() *BranchResponse {
	return &BranchResponse{
		ID:      b.ID,
		Name:    b.Name,
		Code:    b.Code,
		Address: b.Address,
		Phone:   b.Phone,
		Email:   b.Email,
	}
}

