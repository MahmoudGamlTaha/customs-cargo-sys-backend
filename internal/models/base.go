package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        int64     `json:"id" db:"id" gorm:"type:bigint;primary_key;"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

type BaseModelDto struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type NullString sql.NullString

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON allows NullString to accept a JSON string or null
func (ns *NullString) UnmarshalJSON(data []byte) error {
	// Accept null
	if string(data) == "null" {
		*ns = NullString(sql.NullString{String: "", Valid: false})
		return nil
	}
	// Accept plain JSON string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ns = NullString(sql.NullString{String: s, Valid: true})
	return nil
}

// Scan implements the sql.Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*ns = NullString(s)
	return nil
}

// Value implements the driver.Valuer interface for NullString
func (ns NullString) Value() (driver.Value, error) {
	s := sql.NullString(ns)
	return s.Value()
}

// NullInt64 wraps sql.NullInt64 with JSON marshal/unmarshal support
type NullInt64 sql.NullInt64

// MarshalJSON implements json.Marshaler for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON implements json.Unmarshaler for NullInt64
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*ni = NullInt64(sql.NullInt64{Int64: 0, Valid: false})
		return nil
	}
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*ni = NullInt64(sql.NullInt64{Int64: i, Valid: true})
	return nil
}

// Scan implements the sql.Scanner interface for NullInt64
func (ni *NullInt64) Scan(value interface{}) error {
	var n sql.NullInt64
	if err := n.Scan(value); err != nil {
		return err
	}
	*ni = NullInt64(n)
	return nil
}

// Value implements the driver.Valuer interface for NullInt64
func (ni NullInt64) Value() (driver.Value, error) {
	n := sql.NullInt64(ni)
	return n.Value()
}

// UserRole represents the user role enum

const (
	RoleClient      = "client"
	RoleStaff       = "staff"
	RoleAdmin       = "admin"
	RoleAccountant  = "accountant"
	RoleAuditor     = "auditor"
	RoleBranchAdmin = "branch_admin"
)

// RequestStatus represents the request status enum
type RequestStatus string

const (
	StatusSubmitted RequestStatus = "ISSUED"
	StatusApproved  RequestStatus = "APPROVED"
	StatusRejected  RequestStatus = "REJECTED"
	StatusPaid      RequestStatus = "PAID"
)

// IsValidRole checks if a role string is valid
func IsValidRole(role string) bool {
	switch role {
	case RoleClient, RoleStaff, RoleAdmin, RoleAccountant, RoleAuditor, RoleBranchAdmin:
		return true
	default:
		return false
	}
}

// IsValid checks if the request status is valid
func (s RequestStatus) IsValid() bool {
	switch s {
	case StatusSubmitted, StatusApproved, StatusRejected, StatusPaid:
		return true
	default:
		return false
	}
}

// String returns the string representation of UserRole

// String returns the string representation of RequestStatus
func (s RequestStatus) String() string {
	return string(s)
}
