package models

import (
	"time"
)

// RequestDetail represents the detailed information for a request
// Note: client_id and user_id both reference users table
// Relationships are optional and populated when needed
type RequestDetail struct {
	BaseModel

	ClientName       NullString `json:"client_name" db:"client_name"`
	TransferDetail   NullString `json:"transfer_detail" db:"transfer_detail"`
	Signs            NullString `json:"signs" db:"signs"`
	NumberOfParcel   *int       `json:"number_of_parcel" db:"number_of_parcel"`
	RequestID        int64      `json:"request_id" db:"request_id"`
	Description      NullString `json:"description" db:"description"`
	Weight           *float64   `json:"weight" db:"weight"`
	NetWeight        *float64   `json:"net_weight" db:"net_weight"`
	InvoiceNumber    NullString `json:"invoice_number" db:"invoice_number"`
	InvoiceDate      *time.Time `json:"invoice_date" db:"invoice_date"`
	ClientID         *int64     `json:"client_id" db:"client_id"`
	UserID           *int64     `json:"user_id" db:"user_id"`
	CompanyName      NullString `json:"company_name" db:"company_name"`
	CompanyNameEn    NullString `json:"company_name_en" db:"company_name_en"`
	CommercialNumber NullString `json:"commerical_number" db:"commerical_number"`
	ActivityType     NullString `json:"activity_type" db:"activity_type"`
	Address          NullString `json:"address" db:"address"`
	PhoneNumber      NullString `json:"phone_number" db:"phone_number"`
	Email            NullString `json:"email" db:"email"`
	IdentityNumber   NullString `json:"identity_number" db:"identity_number"`
	MobileNumber     NullString `json:"mobile_number" db:"mobile_number"`
	Quantity         *float64   `json:"quantity" db:"quantity"`
	Value            *float64   `json:"item_cost" db:"item_cost"`
	ForOfficialUse   NullString `json:"for_official_use" db:"for_official_use"`
	CountryProducer  NullString `json:"country_producer" db:"country_producer"`
	StandardOfOrigin NullString `json:"standard_of_origin" db:"standard_of_origin"`
	Extra            NullString `json:"extra" db:"extra"`
	ForignItemCost   NullString `json:"foreign_items_cost" db:"foreign_items_cost"`

	// Relationships (populated when needed)
	Request *Request `json:"request,omitempty"`
	Client  *User    `json:"client,omitempty"`
	User    *User    `json:"user,omitempty"`
}

// RequestDetailDto represents the response payload for request detail data
type RequestDetailDto struct {
	ID               *int64        `json:"id"`
	Title            NullString    `json:"title"`
	ClientName       NullString    `json:"client_name"`
	TransferDetail   NullString    `json:"transfer_detail"`
	Signs            NullString    `json:"signs"`
	NumberOfParcel   *int          `json:"number_of_parcel"`
	RequestID        int64         `json:"request_id"`
	Description      NullString    `json:"description"`
	Weight           *float64      `json:"weight"`
	ClientID         NullInt64     `json:"client_id"`
	UserID           NullInt64     `json:"user_id"`
	Quantity         *float64      `json:"quantity"`
	NetWeight        *float64      `json:"net_weight"`
	InvoiceNumber    NullString    `json:"invoice_number"`
	InvoiceDate      *time.Time    `json:"invoice_date"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	Client           *UserResponse `json:"client,omitempty"`
	User             *UserResponse `json:"user,omitempty"`
	CompanyName      NullString    `json:"company_name,omitempty"`
	CommercialNumber NullString    `json:"commerical_number,omitempty"`
	ActivityType     NullString    `json:"activity_type,omitempty"`
	Address          NullString    `json:"address,omitempty"`
	PhoneNumber      NullString    `json:"phone_number,omitempty"`
	Email            NullString    `json:"email,omitempty"`
	IdentityNumber   NullString    `json:"identity_number,omitempty"`
	MobileNumber     NullString    `json:"mobile_number,omitempty"`
	CompanyNameEn    NullString    `json:"company_name_en,omitempty"`
	QRIdentifier     *string       `json:"qr_identifier,omitempty"`
	Value            *float64      `json:"item_cost,omitempty"`
	CountryProducer  NullString    `json:"country_producer,omitempty"`
	RequestTypeName  NullString    `json:"request_type_name,omitempty"`
	ForOfficialUse   NullString    `json:"for_official_use,omitempty"`
	StandardOfOrigin NullString    `json:"standard_of_origin,omitempty"`
	ExporterName     NullString    `json:"exporter_name,omitempty"`
	SerialNumber     NullString    `json:"serial_number,omitempty"`
	Extra            NullString    `json:"extra,omitempty"`
	ForignItemCost   NullString    `json:"foreign_items_cost,omitempty"`
}

// ToResponse converts RequestDetail to RequestDetailDto
func (rd *RequestDetailDto) ToDto() *RequestDetailDto {
	resp := &RequestDetailDto{
		ID:               rd.ID,
		ClientName:       rd.ClientName,
		TransferDetail:   rd.TransferDetail,
		Signs:            rd.Signs,
		NumberOfParcel:   rd.NumberOfParcel,
		RequestID:        rd.RequestID,
		Description:      rd.Description,
		Weight:           rd.Weight,
		NetWeight:        rd.NetWeight,
		InvoiceNumber:    rd.InvoiceNumber,
		InvoiceDate:      rd.InvoiceDate,
		ClientID:         rd.ClientID,
		UserID:           rd.UserID,
		CreatedAt:        rd.CreatedAt,
		UpdatedAt:        rd.UpdatedAt,
		CompanyName:      rd.CompanyName,
		CommercialNumber: rd.CommercialNumber,
		ActivityType:     rd.ActivityType,
		Address:          rd.Address,
		PhoneNumber:      rd.PhoneNumber,
		Email:            rd.Email,
		IdentityNumber:   rd.IdentityNumber,
		MobileNumber:     rd.MobileNumber,
		CompanyNameEn:    rd.CompanyNameEn,
		CountryProducer:  rd.CountryProducer,
		StandardOfOrigin: rd.StandardOfOrigin,
		Extra:            rd.Extra,
		ForignItemCost:   rd.ForignItemCost,
	}
	return resp
}
