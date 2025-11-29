package service

import (
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
	"Chumber-Workflow-System/pkg/client"
	"Chumber-Workflow-System/pkg/crypto"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// RequestService handles request business logic
type RequestService struct {
	requestRepo     *repository.RequestRepository
	userRepo        *repository.UserRepository
	requestTypeRepo *repository.RequestTypeRepository
	activityRepo    *repository.UserActivityRepository
	branchRepo      *repository.BranchRepository
	requestClient   *client.RequestClient
}

// NewRequestService creates a new request service
func NewRequestService(
	requestRepo *repository.RequestRepository,
	userRepo *repository.UserRepository,
	requestTypeRepo *repository.RequestTypeRepository,
	activityRepo *repository.UserActivityRepository,
	branchRepo *repository.BranchRepository,
) *RequestService {
	return &RequestService{
		requestRepo:     requestRepo,
		userRepo:        userRepo,
		requestTypeRepo: requestTypeRepo,
		activityRepo:    activityRepo,
		branchRepo:      branchRepo,
		requestClient:   client.NewRequestClient(),
	}
}

// CreateRequest creates a new request
func (s *RequestService) CreateRequest(req *models.CreateRequestRequest, currentUserID int64, currentUserRole string) (*models.Request, error) {
	// Get current user

	currentUser, err := s.userRepo.GetByID(currentUserID)
	if err != nil {
		return nil, err
	}
	req.BranchID = *currentUser.BranchID
	req.UserID = currentUserID
	var clientID int64

	if currentUserRole == models.RoleClient {
		// Client creating their own request
		clientID = currentUserID
		if req.ClientID <= 0 {
			return nil, errors.New("client_id is required when staff creates request")

		}
	} else if currentUserRole == models.RoleStaff || currentUserRole == models.RoleAdmin || currentUserRole == models.RoleBranchAdmin {
		// Staff/Admin creating request for a client

		clientID = req.ClientID

		// Verify client exists and is a client
		/*	client, err := s.userRepo.GetByID(clientID)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, errors.New("client not found")
				}
				return nil, err
			}*/

		/*	if client.Role != models.RoleClient {
			return nil, errors.New("specified user is not a client")
		}*/

		//clientID = client.ID
	} else {
		return nil, errors.New("unauthorized to create requests")
	}
	fmt.Println("currentUser.BranchID8:", currentUser.BranchID)
	fmt.Println("currentUser.BranchID:", req.BranchID)
	// Create request
	transaction := s.requestRepo.BeginTransaction()
	request := &models.Request{
		Title:           req.Title,
		Description:     req.Description,
		Status:          models.StatusSubmitted,
		ClientID:        &clientID,
		CreatedByUserID: &currentUserID,
		RequestTypeID:   &req.RequestTypeID, // 4
		BranchID:        currentUser.BranchID,
		ExporterName:    &req.ExporterName,
	}

	err = s.requestRepo.Create(request)
	if err != nil {
		return nil, err
	}

	// Ensure at least one request detail is provided
	if len(req.RequestDetails) == 0 {
		return nil, errors.New("request_details is required")
	}

	// Create detail rows for the request
	for _, d := range req.RequestDetails {
		var numParcelPtr *int
		if d.NumberOfParcel != nil && *d.NumberOfParcel != 0 {
			numParcelPtr = d.NumberOfParcel
		}
		var weightPtr *float64
		if d.Weight != nil && *d.Weight != 0 {
			weightPtr = d.Weight
		}
		var netWeightPtr *float64
		if d.NetWeight != nil && *d.NetWeight != 0 {
			netWeightPtr = d.NetWeight
		}

		rd := &models.RequestDetail{
			RequestID:        request.ID,
			ClientID:         &clientID,
			UserID:           &currentUserID,
			ClientName:       d.ClientName,
			TransferDetail:   d.TransferDetail,
			Signs:            d.Signs,
			NumberOfParcel:   numParcelPtr,
			Description:      d.Description,
			Weight:           weightPtr,
			NetWeight:        netWeightPtr,
			InvoiceNumber:    d.InvoiceNumber,
			InvoiceDate:      d.InvoiceDate,
			CompanyName:      d.CompanyName,
			CompanyNameEn:    d.CompanyNameEn,
			CommercialNumber: d.CommercialNumber,
			ActivityType:     d.ActivityType,
			Address:          d.Address,
			PhoneNumber:      d.PhoneNumber,
			Email:            d.Email,
			IdentityNumber:   d.IdentityNumber,
			MobileNumber:     d.MobileNumber,
			Quantity:         d.Quantity,
			Value:            d.Value,
			CountryProducer:  d.CountryProducer,
			StandardOfOrigin: d.StandardOfOrigin,
			ForOfficialUse:   d.ForOfficialUse,
			Extra:            d.Extra,
			ForignItemCost:   d.ForignItemCost,
		}

		if err := s.requestRepo.CreateRequestDetail(rd, request.ID); err != nil {
			return nil, err
		}
	}
	transaction = s.requestRepo.CommitTransaction(transaction)
	if transaction == nil {
		return nil, errors.New("failed to commit transaction")
	}
	return request, nil
}

// GetRequest retrieves a request by ID with access control
func (s *RequestService) GetRequest(id int64, currentUserID int64, currentUserRole string) (*models.Request, error) {
	fmt.Println("id:", id)

	// Get user's branch ID for branch-based filtering
	var userBranchID string
	if currentUserRole != models.RoleAdmin {
		user, err := s.userRepo.GetByID(currentUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %w", err)
		}

		if user.BranchID != nil {
			userBranchID = fmt.Sprintf("%d", *user.BranchID)
		}
	}

	// Use the repository method with role-based filtering
	request, err := s.requestRepo.GetByID(id, currentUserRole, userBranchID)
	if err != nil {
		fmt.Println("err:", err)
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	// Check access permissions
	if !s.canAccessRequest(request, currentUserID, currentUserRole) {
		return nil, errors.New("access denied")
	}

	return request, nil
}

// GetRequestWithDetails retrieves a request with user details
func (s *RequestService) GetRequestWithDetails(id int64, currentUserID int64, currentUserRole string, secret string) (*models.RequestDetailDto, error) {
	// Get user's branch ID for branch-based filtering
	var userBranchID string
	if currentUserRole != models.RoleAdmin {
		user, err := s.userRepo.GetByID(currentUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %w", err)
		}

		if user.BranchID != nil {
			userBranchID = fmt.Sprintf("%d", *user.BranchID)
		}
	}

	// First check if user can access this request
	_, err := s.requestRepo.GetByID(id, currentUserRole, userBranchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	// Get detailed request
	requestDetails, err := s.requestRepo.GetWithDetails(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request details not found")
		}
		return nil, err
	}

	// The QR identifier should already be stored in the database
	// If not present, generate it using the secret parameter
	if requestDetails.QRIdentifier == nil && secret != "" {
		encryptedID, err := crypto.EncryptID(id, secret)
		if err == nil {
			requestDetails.QRIdentifier = &encryptedID
		}
	}

	return requestDetails, nil
}

// GetRequestWithDetails retrieves a request with user details
func (s *RequestService) GetRequestWithDetailsUsingIdentifier(encryptedID string, secret string) (*models.RequestDetailDto, error) {
	// Decrypt the identifier to get the actual request ID
	id, err := crypto.DecryptID(encryptedID, secret)
	if err != nil {
		return nil, errors.New("invalid identifier")
	}

	// First check if user can access this request
	req, err := s.requestRepo.GetByID(id, "", "")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	if req.Status != models.RequestStatus(models.StatusPaid.String()) {
		return nil, errors.New("request not found")
	}

	// Get detailed request
	requestDetails, err := s.requestRepo.GetWithDetails(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request details not found")
		}
		return nil, err
	}

	return requestDetails, nil
}

// UpdateRequest updates a request (only if pending and by client or staff)
func (s *RequestService) UpdateRequest(id int64, req *models.UpdateRequestRequest, currentUserID int64, currentUserRole string) (*models.Request, error) {
	// Get user's branch ID for branch-based filtering
	var userBranchID string
	if currentUserRole != models.RoleAdmin {
		user, err := s.userRepo.GetByID(currentUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %w", err)
		}

		if user.BranchID != nil {
			userBranchID = fmt.Sprintf("%d", *user.BranchID)
		}
	}

	// Get existing request with branch filtering
	request, err := s.requestRepo.GetByID(id, currentUserRole, userBranchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	// Check if request can be edited
	if !request.CanBeEdited() {
		return nil, errors.New("request cannot be edited after approval or rejection")
	}

	// Check access permissions
	if !s.canEditRequest(request, currentUserID, currentUserRole) {
		return nil, errors.New("access denied")
	}

	// Update fields if provided
	if req.Title != "" {
		request.Title = req.Title
	}
	if req.Description != "" {
		request.Description = req.Description
	}

	// Update request
	err = s.requestRepo.Update(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// ApproveRequest approves a request (staff/admin only)
func (s *RequestService) ApproveRequest(id int64, currentUserID int64, currentUserRole string) (*models.Request, error) {
	// Only staff and admin can approve requests
	// Get user's branch ID for branch-based filtering
	var userBranchID string
	if currentUserRole != models.RoleAdmin {
		user, err := s.userRepo.GetByID(currentUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %w", err)
		}

		if user.BranchID != nil {
			userBranchID = fmt.Sprintf("%d", *user.BranchID)
		}
	}

	// Get request with branch filtering
	request, err := s.requestRepo.GetByID(id, currentUserRole, userBranchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	// Check if request can be approved
	if !request.CanBeApproved() {
		return nil, errors.New("request cannot be approved")
	}

	// Verify the staff member exists and is staff/admin
	staff, err := s.userRepo.GetByID(currentUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("staff member not found")
		}
		return nil, err
	}

	if staff.Role == nil || (*staff.Role != models.RoleStaff && *staff.Role != models.RoleAdmin && *staff.Role != models.RoleAuditor && *staff.Role != models.RoleBranchAdmin) {
		return nil, errors.New("approver must be staff or admin or branch admin")
	}
	// Approve request (this will trigger serial number generation in database)
	err = s.requestRepo.Approve(id, currentUserID)
	if err != nil {
		return nil, err
	}

	// Get updated request - use admin role to bypass branch filtering for result
	updatedRequest, err := s.requestRepo.GetByID(id, models.RoleAdmin, "")
	if err != nil {
		return nil, err
	}

	return updatedRequest, nil
}

// RejectRequest rejects a request (staff/admin only)
func (s *RequestService) RejectRequest(id int64, currentUserID int64, currentUserRole string, rejectionReason string) (*models.Request, error) {
	// Only staff and admin can reject request

	// Get user's branch ID for branch-based filtering
	var userBranchID string
	if currentUserRole != models.RoleAdmin {
		user, err := s.userRepo.GetByID(currentUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %w", err)
		}

		if user.BranchID != nil {
			userBranchID = fmt.Sprintf("%d", *user.BranchID)
		}
	}

	// Get request with branch filtering
	request, err := s.requestRepo.GetByID(id, currentUserRole, userBranchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	// Check if request can be rejected
	if !request.CanBeRejected() {
		return nil, errors.New("request cannot be rejected")
	}

	// Parse approved by staff ID (who is doing the rejection)
	approvedByStaffID := currentUserID

	// Verify the staff member exists and is staff/admin
	staff, err := s.userRepo.GetByID(approvedByStaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("staff member not found")
		}
		return nil, err
	}

	if staff.Role == nil || (*staff.Role != models.RoleStaff && *staff.Role != models.RoleAdmin && *staff.Role != models.RoleAuditor && *staff.Role != models.RoleBranchAdmin && *staff.Role != models.RoleAccountant) {
		return nil, errors.New("rejector must be staff, admin, auditor, branch admin or accountant")
	}

	// Reject request
	err = s.requestRepo.Reject(id, approvedByStaffID, rejectionReason)
	if err != nil {
		return nil, err
	}

	// Get updated request - use admin role to bypass branch filtering for result
	updatedRequest, err := s.requestRepo.GetByID(id, models.RoleAdmin, "")
	if err != nil {
		return nil, err
	}
	return updatedRequest, nil
}

func (s *RequestService) MarkAsPaid(id int64, currentUserID int64, currentUserRole string, secret string) (*models.Request, error) {
	// Only staff and admin can mark request as paid

	// Get user's branch ID for branch-based filtering
	var userBranchID int64
	if currentUserRole != models.RoleAdmin || currentUserRole != models.RoleBranchAdmin {
		user, err := s.userRepo.GetByID(currentUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %w", err)
		}

		if user.BranchID != nil {
			userBranchID = *user.BranchID
		}
	}

	// Get request with branch filtering (just to verify it exists and user has access)
	request, err := s.requestRepo.GetByID(id, currentUserRole, strconv.FormatInt(userBranchID, 10))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	// Get request type to calculate final amount (price * ratio)
	if request.RequestTypeID == nil {
		return nil, errors.New("request has no associated request type")
	}

	requestType, err := s.requestTypeRepo.GetByID(*request.RequestTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get request type: %w", err)
	}

	// Calculate final amount: price * ratio
	calculatedAmount := requestType.Price * requestType.Ratio

	// Mark request as paid with calculated amount
	err = s.requestRepo.MarkAsPaid(id, currentUserID, calculatedAmount)
	if err != nil {
		return nil, err
	}
	s.requestRepo.GenerateSerialNumber(id, userBranchID)
	// Get updated request - use admin role to bypass branch filtering for result
	updatedRequest, err := s.requestRepo.GetByID(id, models.RoleAdmin, "")
	if err != nil {
		return nil, err
	}

	// Generate and store encrypted QR identifier after request is created
	if secret != "" {
		go func() {
			encryptedID, err := crypto.EncryptID(updatedRequest.ID, secret)
			if err == nil {
				updatedRequest.QRIdentifier = &encryptedID
				// Update the request with the QR identifier
				err = s.requestRepo.UpdateIdentifier(updatedRequest)
			}
			if err != nil {
				// Log error but don't fail the request creation
				fmt.Printf("Failed to update QR identifier: %v\n", err)
			}
		}()
	}

	return updatedRequest, nil
}

// ListRequests retrieves requests with pagination and filtering
func (s *RequestService) ListRequests(params models.FilterParams, pagination models.PaginationParams, currentUserID int64, currentUserRole string) ([]*models.RequestResponse, *models.PaginationResponse, error) {
	var requests []*models.RequestResponse
	var total int
	var err error

	if currentUserRole == models.RoleClient {
		// Clients can only see their own requests
		requests, total, err = s.requestRepo.ListByClientID(currentUserID, params, pagination, currentUserRole)
	} else {
		// Staff and admin can see all requests
		requests, total, err = s.requestRepo.List(params, pagination, currentUserRole)
	}

	if err != nil {
		return nil, nil, err
	}

	paginationResp := models.CalculatePagination(pagination.Page, pagination.PageSize, total)
	return requests, &paginationResp, nil
}

func (s *RequestService) ListRequestsCount(params models.FilterParams, currentUserID int64, currentUserRole string) (map[string]int, error) {
	var total int
	var err error

	if currentUserRole == models.RoleClient {
		// Clients can only see their own requests
		total, err = s.requestRepo.ListDetailsCount(currentUserID, params, currentUserRole)
	} else {
		// Staff and admin can see all requests
		total, err = s.requestRepo.ListDetailsCount(currentUserID, params, currentUserRole)
	}

	if err != nil {
		return nil, err
	}

	return map[string]int{"total": total}, nil
}

// GetRequestHistory retrieves request history by request ID
func (s *RequestService) GetRequestHistory(requestID int64, currentUserID int64, currentUserRole string) ([]*models.RequestHistory, error) {
	// Get user's branch ID for branch-based filtering
	var userBranchID string
	if currentUserRole != models.RoleAdmin {
		user, err := s.userRepo.GetByID(currentUserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %w", err)
		}

		if user.BranchID != nil {
			userBranchID = fmt.Sprintf("%d", *user.BranchID)
		}
	}

	// First check if request exists and user has access
	_, err := s.requestRepo.GetByID(requestID, currentUserRole, userBranchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	// Get request history
	history, err := s.requestRepo.GetHistory(requestID)
	if err != nil {
		return nil, err
	}

	return history, nil
}

// DeleteRequest deletes a request (admin only, and only if pending)
func (s *RequestService) DeleteRequest(id int64, currentUserID int64, currentUserRole string) error {
	// Only admin can delete requests
	if currentUserRole != models.RoleAdmin {
		return errors.New("only admin can delete requests")
	}

	// Get request - since only admins can delete, use admin role and empty branch ID to bypass filtering
	request, err := s.requestRepo.GetByID(id, models.RoleAdmin, "")
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("request not found")
		}
		return err
	}

	// Only allow deletion of pending requests
	if request.Status != models.StatusSubmitted {
		return errors.New("only Issued requests can be deleted")
	}

	return s.requestRepo.Delete(id)
}

// canAccessRequest checks if a user can access a specific request
func (s *RequestService) canAccessRequest(request *models.Request, userID int64, userRole string) bool {
	switch userRole {
	case models.RoleClient:
		// Clients can only access their own requests
		return *request.ClientID == userID
	case models.RoleStaff, models.RoleAdmin:
		// Staff and admin can access all requests
		return true
	default:
		return false
	}
}

// canEditRequest checks if a user can edit a specific request
func (s *RequestService) canEditRequest(request *models.Request, userID int64, userRole string) bool {
	switch userRole {
	case models.RoleClient:
		// Clients can only edit their own pending requests
		return *request.ClientID == userID && request.Status == models.StatusSubmitted
	case models.RoleStaff, models.RoleAdmin:
		// Staff and admin can edit pending requests
		return request.Status == models.StatusSubmitted
	default:
		return false
	}
}

// GetRequestBySerial retrieves a request by its serial number from external API
func (s *RequestService) GetRequestBySerial(serialNumber string, userID int64, username string, branchID int64) (*models.Request, error) {
	// Fetch request from external API
	request, err := s.requestClient.GetRequestBySerial(serialNumber)
	if err != nil {
		return nil, err
	}

	return request, nil
}
