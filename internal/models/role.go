package models

// Role represents a role entity
type Role struct {
	BaseModel
	NameAR      string        `json:"name_ar" db:"name_ar"`
	NameEN      string        `json:"name_en" db:"name_en"`
	Code        string        `json:"code" db:"code"`
	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:permission_role;"`
}

// CreateRoleRequest represents the request body for creating a new role
type CreateRoleRequest struct {
	NameAR string `json:"name_ar" binding:"required"`
	NameEN string `json:"name_en" binding:"required"`
	Code   string `json:"code" binding:"required,alpha"`
}
