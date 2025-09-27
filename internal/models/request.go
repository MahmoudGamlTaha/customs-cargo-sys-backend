package models

import (
	"time"
)

// Request represents a request entity
type Request struct {
	BaseModel
	SerialNumber     *string          `json:"serial_number" db:"serial_number"`
	Title            string           `json:"title" db:"title" validate:"required,min=5,max=255"`
	Description      string           `json:"description" db:"description" validate:"required,min=10"`
	Status           RequestStatus    `json:"status" db:"status"`
	ClientID         *int64           `json:"client_id" db:"client_id" validate:"omitempty,gt=0"`
	CreatedByUserID  *int64           `json:"created_by" db:"created_by"`
	ApprovedByUserID *int64           `json:"approved_by" db:"approved_by"`
	RejectionReason  *string          `json:"rejection_reason" db:"rejection_reason"`
	ApprovedAt       *time.Time       `json:"approved_at" db:"approved_at"`
	RejectedAt       *time.Time       `json:"rejected_at" db:"rejected_at"`
	RequestDetailID  []*RequestDetail `json:"request_detail_id" db:"request_detail_id"`
	RequestTypeID    *int64           `json:"request_type_id" db:"request_type_id"`
	BranchID         *int64           `json:"branch_id" db:"branch_id"`
	ExporterName     *string          `json:"exporter_name" db:"exporter_name"`
	RequestTypeName  *string          `json:"request_type_name" db:"request_type_name"`
	QRIdentifier     *string          `json:"qr_identifier" db:"qr_identifier"`
	// Relationships (populated when needed)
	Client         *User            `json:"client,omitempty"`
	CreatedByUser  *User            `json:"created_by_user,omitempty"`
	ApprovedByUser *User            `json:"approved_by_user,omitempty"`
	RequestDetails []*RequestDetail `json:"request_details,omitempty"`
	Ratio          *float64         `json:"ratio,omitempty"`
}

// RequestDetails represents request with user details

// CreateRequestRequest represents the request payload for creating a request
type CreateRequestRequest struct {
	Title          string              `json:"title" validate:"required,min=5,max=255"`
	Description    string              `json:"description" validate:"required,min=10"`
	ClientID       int64               `json:"client_id" validate:"omitempty,gt=0"`
	UserID         int64               `json:"user_id" validate:"omitempty,gt=0"`
	RequestTypeID  int64               `json:"request_type_id" validate:"omitempty,gt=0"`
	ExporterName   string              `json:"exporter_name" validate:"omitempty,min=5,max=255"`
	BranchID       int64               `json:"branch_id" validate:"omitempty,gt=0"`
	RequestDetails []*RequestDetailDto `json:"request_details" validate:"required,dive"`

	// Optional for staff creating for client
}

// UpdateRequestRequest represents the request payload for updating a request
type UpdateRequestRequest struct {
	Title       string `json:"title" validate:"omitempty,min=5,max=255"`
	Description string `json:"description" validate:"omitempty,min=10"`
}

// ApproveRequestRequest represents the request payload for approving a request
type ApproveRequestRequest struct {
	ApprovedByStaffID int64 `json:"approved_by_staff_id" validate:"required,gt=0"`
}

// RejectRequestRequest represents the request payload for rejecting a request
type RejectRequestRequest struct {
	RejectionReason   string `json:"rejection_reason" validate:"required,min=10"`
	ApprovedByStaffID int64  `json:"approved_by_staff_id" validate:"required,gt=0"`
}

// RequestResponse represents the response payload for request data
type RequestResponse struct {
	ID               int64               `json:"id"`
	SerialNumber     *string             `json:"serial_number"`
	Title            string              `json:"title"`
	Description      string              `json:"description"`
	Status           RequestStatus       `json:"status"`
	RejectionReason  *string             `json:"rejection_reason"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	ApprovedAt       *time.Time          `json:"approved_at"`
	RejectedAt       *time.Time          `json:"rejected_at"`
	Client           *UserResponse       `json:"client,omitempty"`
	CreatedByStaff   *UserResponse       `json:"created_by_staff,omitempty"`
	ApprovedByStaff  *UserResponse       `json:"approved_by_staff,omitempty"`
	CreateByUserID   *int64              `json:"create_by_user_id,omitempty"`
	ApprovedByUserID *int64              `json:"approved_by_user_id,omitempty"`
	RequestTypeID    *int64              `json:"request_type_id,omitempty"`
	RequestTypeName  *string             `json:"request_type_name,omitempty"`
	ExporterName     *string             `json:"exporter_name,omitempty"`
	ServiceType      *string             `json:"service_type,omitempty"`
	RequestDetails   []*RequestDetailDto `json:"request_details"`
	QRIdentifier     *string             `json:"qr_identifier"`
	Ratio            *float64            `json:"ratio,omitempty"`
	BranchID         *int64              `json:"branch_id,omitempty"`
	CalculatedAmount *float64            `json:"calculated_amount,omitempty"`
}

// RequestListResponse represents the response payload for request list
type RequestListResponse struct {
	Requests   []*RequestResponse `json:"requests"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// RequestHistory represents a request history entry
type RequestHistory struct {
	BaseModel
	RequestID       int64          `json:"request_id" db:"request_id"`
	OldStatus       *RequestStatus `json:"old_status" db:"old_status"`
	NewStatus       RequestStatus  `json:"new_status" db:"new_status"`
	ChangedByUserID int64          `json:"changed_by_user_id" db:"changed_by_user_id"`
	ChangeReason    *string        `json:"change_reason" db:"change_reason"`

	// Relationships
	Request       *Request `json:"request,omitempty"`
	ChangedByUser *User    `json:"changed_by_user,omitempty"`
}

// ToResponse converts Request to RequestResponse
func (r *Request) ToResponse() *RequestResponse {
	qrIdentifier := ""
	if r.QRIdentifier != nil {
		qrIdentifier = *r.QRIdentifier
	}

	response := &RequestResponse{
		ID:               r.ID,
		SerialNumber:     r.SerialNumber,
		Title:            r.Title,
		Description:      r.Description,
		Status:           r.Status,
		RejectionReason:  r.RejectionReason,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
		ApprovedAt:       r.ApprovedAt,
		RejectedAt:       r.RejectedAt,
		CreateByUserID:   r.CreatedByUserID,
		ApprovedByUserID: r.ApprovedByUserID,
		RequestTypeID:    r.RequestTypeID,
		RequestDetails:   r.toDetailResponse(r.RequestDetails),
		QRIdentifier:     &qrIdentifier,
		Ratio:            r.Ratio,
	}
	return response
}

func (r *Request) toDetailResponse(requestDetails []*RequestDetail) []*RequestDetailDto {
	var requestDetailDto []*RequestDetailDto
	for _, requestDetail := range requestDetails {
		requestDetailDto = append(requestDetailDto, &RequestDetailDto{
			RequestID:        r.ID,
			Signs:            requestDetail.Signs,
			NumberOfParcel:   requestDetail.NumberOfParcel,
			Description:      requestDetail.Description,
			Weight:           requestDetail.Weight,
			NetWeight:        requestDetail.NetWeight,
			InvoiceNumber:    requestDetail.InvoiceNumber,
			InvoiceDate:      requestDetail.InvoiceDate,
			TransferDetail:   requestDetail.TransferDetail,
			ClientName:       requestDetail.ClientName,
			CommercialNumber: requestDetail.CommercialNumber,
			ActivityType:     requestDetail.ActivityType,
			Address:          requestDetail.Address,
			PhoneNumber:      requestDetail.PhoneNumber,
			Email:            requestDetail.Email,
			IdentityNumber:   requestDetail.IdentityNumber,
			MobileNumber:     requestDetail.MobileNumber,
			CompanyName:      requestDetail.CompanyName,
		})
	}
	return requestDetailDto
}

// IsPending checks if the request is pending
func (r *Request) IsPending() bool {
	return r.Status == StatusSubmitted
}

// IsApproved checks if the request is approved
func (r *Request) IsApproved() bool {
	return r.Status == StatusApproved
}

// IsRejected checks if the request is rejected
func (r *Request) IsRejected() bool {
	return r.Status == StatusRejected
}

// CanBeEdited checks if the request can be edited
func (r *Request) CanBeEdited() bool {
	return r.Status == StatusSubmitted
}

// CanBeApproved checks if the request can be approved
func (r *Request) CanBeApproved() bool {
	return r.Status == StatusSubmitted
}

// CanBeRejected checks if the request can be rejected
func (r *Request) CanBeRejected() bool {
	return r.Status == StatusSubmitted
}
