package models

import (
	"time"
)

// UserActivities represents user activities
type UserActivities struct {
	BaseModel
	Description  string     `json:"description" db:"description"`
	Action       string     `json:"action" db:"action"`
	UserId       int64      `json:"user_id" db:"user_id"`
	Username     string     `json:"username" db:"username"`
	Module       string     `json:"module" db:"module"`
	EntityId     *int64     `json:"entity_id" db:"entity_id"`
	EntityName   string     `json:"entity_name" db:"entity_name"`
	SerialNumber string     `json:"serial_number" db:"serial"`
	BranchID     *int64     `json:"branch_id" db:"branch_id"`
	CreatedAt    *time.Time `json:"created_at" db:"created_at"`
}

// UserActivitiesResponse represents the response payload for user activities data
type UserActivitiesResponse struct {
	ID           int64  `json:"id"`
	Description  string `json:"description"`
	Action       string `json:"action,omitempty"`
	UserId       int64  `json:"user_id"`
	Module       string `json:"module"`
	Username     string `json:"username"`
	EntityId     *int64 `json:"entity_id"`
	EntityName   string `json:"entity_name"`
	SerialNumber string `json:"serial_number"`
	BranchID     *int64 `json:"branch_id"`
}

// ToResponse converts UserActivities to UserActivitiesResponse
func (rt *UserActivities) ToResponse() *UserActivitiesResponse {
	return &UserActivitiesResponse{
		ID:           rt.ID,
		Description:  rt.Description,
		Action:       rt.Action,
		UserId:       rt.UserId,
		Module:       rt.Module,
		Username:     rt.Username,
		EntityId:     rt.EntityId,
		EntityName:   rt.EntityName,
		SerialNumber: rt.SerialNumber,
		BranchID:     rt.BranchID,
	}
}
