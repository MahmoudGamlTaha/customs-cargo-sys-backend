package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Date is a custom type for handling date-only values (without time)
type Date time.Time

func (d Date) ToTime() time.Time {
	return time.Time(d)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	// Parse the date string in the format 2006-01-02
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return json.Marshal(t.Format("2006-01-02"))
}

// String returns the date as a string in YYYY-MM-DD format
func (d Date) String() string {
	t := time.Time(d)
	return t.Format("2006-01-02")
}

// Value implements the driver.Valuer interface
func (d Date) Value() (driver.Value, error) {
	t := time.Time(d)
	return t, nil
}

// Scan implements the sql.Scanner interface
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = Date(v)
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		*d = Date(t)
		return nil
	default:
		return fmt.Errorf("unsupported type for Date: %T", value)
	}
}

type Memberships struct {
	BaseModel
	Description    *string       `json:"description" db:"description"`
	Status         RequestStatus `json:"status" db:"status"`
	CreatedBy      int64         `json:"created_by" db:"created_by"`
	BranchId       *int64        `json:"branch_id" db:"branch_id"`
	BranchName     NullString    `json:"branch_name" db:"branch_name"`
	AppDate        *Date         `json:"app_date" db:"app_date"`
	AppNumber      NullString    `json:"app_number" db:"app_number"`
	CancelReason   NullString    `json:"cancel_reason" db:"cancel_reason"`
	RejectReason   NullString    `json:"reject_reason" db:"reject_reason"`
	SerialNumber   NullString    `json:"serial_number" db:"serial_number"`
	CreateDate     *Date         `json:"create_date" db:"create_date"`
	ExpirationDate *Date         `json:"expiration_date" db:"expiration_date"`
}

type MembershipDetails struct {
	BaseModel
	ExporterName   *string `json:"exporter_name" db:"exporter_name"`
	CEOName        *string `json:"ceo_name" db:"ceo_name"`
	MembershipId   int64   `json:"membership_id" db:"membership_id"`
	GeneralManager *string `json:"general_manager" db:"general_manager"`
	Address        *string `json:"address" db:"address"`
	GovShape       *string `json:"gov_shape" db:"gov_shape"`
	TypeName       *string `json:"type_name" db:"type_name"`
}

type MembershipDetailsDto struct {
	BaseModelDto
	ExporterName   *string `json:"exporter_name"`
	CEOName        *string `json:"ceo_name"`
	MembershipId   int64   `json:"membership_id"`
	GeneralManager *string `json:"general_manager"`
	Address        *string `json:"address"`
	GovShape       *string `json:"gov_shape"`
	TypeName       *string `json:"type_name"`
}

type MembershipsDTO struct {
	BaseModelDto
	Description    *string                 `json:"description"`
	Status         RequestStatus           `json:"status"`
	CreatedBy      int64                   `json:"created_by"`
	BranchId       *int64                  `json:"branch_id"`
	BranchName     NullString              `json:"branch_name"`
	AppDate        *Date                   `json:"app_date"`
	AppNumber      NullString              `json:"app_number"`
	CancelReason   NullString              `json:"cancel_reason"`
	RejectReason   NullString              `json:"reject_reason"`
	SerialNumber   NullString              `json:"serial_number"`
	CreateDate     *Date                   `json:"create_date"`
	ExpirationDate *Date                   `json:"expiration_date"`
	IsExpired      bool                    `json:"is_expired"`
	Details        []*MembershipDetailsDto `json:"details"`
}

type RenewMembership struct {
	ExpirationDate *Date `json:"expiration_date" binding:"required"`
}
