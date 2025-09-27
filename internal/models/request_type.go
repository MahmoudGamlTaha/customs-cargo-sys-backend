package models

import (
	"time"
)

// RequestType represents a type/category of request with pricing
type RequestType struct {
	BaseModel
	Name      string     `json:"name" db:"name"`
	NameEn    NullString `json:"name_en" db:"name_en"`
	Price     float64    `json:"price" db:"price"`
	Ratio     float64    `json:"ratio" db:"ratio"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// RequestTypeResponse represents the response payload for request type data
type RequestTypeResponse struct {
	ID     int64      `json:"id"`
	Name   string     `json:"name"`
	NameEn NullString `json:"name_en,omitempty"`
	Price  float64    `json:"price"`
	Ratio  float64    `json:"ratio"`
}

// ToResponse converts RequestType to RequestTypeResponse
func (rt *RequestType) ToResponse() *RequestTypeResponse {
	return &RequestTypeResponse{
		ID:     rt.ID,
		Name:   rt.Name,
		NameEn: rt.NameEn,
		Price:  rt.Price,
		Ratio:  rt.Ratio,
	}
}
