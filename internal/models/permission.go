package models

// Permission represents a permission entity
type Permission struct {
	BaseModel
	NameAR string  `json:"name_ar" db:"name_ar"`
	NameEN string  `json:"name_en" db:"name_en"`
	Code   string  `json:"code" db:"code"`
	Roles  []*Role `json:"roles,omitempty" gorm:"many2many:permission_role;"`
}
