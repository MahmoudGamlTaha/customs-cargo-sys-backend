package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/database"
)

// RequestRepository handles request data operations
type RequestRepository struct {
	db *database.DB
}

func (r *RequestRepository) CommitTransaction(transaction any) any {
	tx := transaction.(*sql.Tx)
	err := tx.Commit()
	if err != nil {
		return nil
	}
	return tx
}

func (r *RequestRepository) BeginTransaction() any {
	tx, err := r.db.Begin()
	if err != nil {
		return nil
	}
	return tx
}

// NewRequestRepository creates a new request repository
func NewRequestRepository(db *database.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

// Create creates a new request
func (r *RequestRepository) Create(request *models.Request) error {
	query := `
		INSERT INTO requests (title, description, status, client_id, created_by, branch_id, exporter_name, request_type_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	var createdBy sql.NullInt64
	if request.CreatedByUserID != nil {
		createdBy = sql.NullInt64{Int64: *request.CreatedByUserID, Valid: true}
	}

	var branchID sql.NullInt64
	if request.BranchID != nil {
		branchID = sql.NullInt64{Int64: *request.BranchID, Valid: true}
	}

	var exporterName sql.NullString
	if request.ExporterName != nil {
		exporterName = sql.NullString{String: *request.ExporterName, Valid: true}
	}

	var requestTypeID sql.NullInt64
	if request.RequestTypeID != nil {
		requestTypeID = sql.NullInt64{Int64: *request.RequestTypeID, Valid: true}
	}

	err := r.db.QueryRow(
		query,
		request.Title,
		request.Description,
		request.Status,
		request.ClientID,
		createdBy,
		branchID,
		exporterName,
		requestTypeID,
	).Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)

	return err
}

// GetByID retrieves a request by ID
func (r *RequestRepository) GetByID(id int64, userRole string, userBranchID string) (*models.Request, error) {
	request := &models.Request{}

	// Base query
	baseQuery := `
		SELECT id, serial_number, title, description, status, client_id, created_by,
		       approved_by, rejection_reason, created_at, updated_at, approved_at, rejected_at,
		       branch_id, exporter_name, request_type_id
		FROM requests 
		WHERE id = $1`

	// Add branch filtering for non-admin users
	query := baseQuery
	args := []interface{}{id}

	if userRole != models.RoleAdmin && userBranchID != "" {
		query += " AND branch_id = $2"
		argNum, _ := strconv.ParseInt(userBranchID, 10, 64)
		args = append(args, argNum)
	}

	var serial models.NullString
	var createdBy models.NullInt64
	var approvedBy models.NullInt64
	var rejection models.NullString
	var approvedAt sql.NullTime
	var rejectedAt sql.NullTime

	var branchID models.NullInt64
	var requestTypeID models.NullInt64

	err := r.db.QueryRow(query, args...).Scan(
		&request.ID,
		&serial,
		&request.Title,
		&request.Description,
		&request.Status,
		&request.ClientID,
		&createdBy,
		&approvedBy,
		&rejection,
		&request.CreatedAt,
		&request.UpdatedAt,
		&approvedAt,
		&rejectedAt,
		&branchID,
		&request.ExporterName,
		&requestTypeID,
	)

	if err != nil {
		return nil, err
	}

	// Set values from NullString/NullInt64 to the struct fields
	if serial.Valid {
		request.SerialNumber = &serial.String
	}
	if createdBy.Valid {
		request.CreatedByUserID = &createdBy.Int64
	}
	if approvedBy.Valid {
		request.ApprovedByUserID = &approvedBy.Int64
	}
	if rejection.Valid {
		request.RejectionReason = &rejection.String
	}
	if approvedAt.Valid {
		request.ApprovedAt = &approvedAt.Time
	}
	if rejectedAt.Valid {
		request.RejectedAt = &rejectedAt.Time
	}

	// Set the new fields
	if branchID.Valid {
		request.BranchID = &branchID.Int64
	}

	if requestTypeID.Valid {
		request.RequestTypeID = &requestTypeID.Int64
	}

	fmt.Println("request", request.ExporterName)

	return request, nil
}

// GetWithDetails retrieves a request with user details
func (r *RequestRepository) GetWithDetails(id int64) (*models.RequestDetailDto, error) {
	requestDetails := &models.RequestDetailDto{}
	query := `
		SELECT 
			r.id as request_id, r.serial_number, r.title, r.description, r.status, r.client_id,
			r.created_by,r.request_type_id, r.created_at, r.updated_at,
			rd.id as detail_id, rd.client_name, rd.transfer_detail, rd.signs, 
			rd.number_of_parcel, rd.description as detail_description, rd.weight, rd.net_weight, 
			rd.invoice_number, rd.invoice_date, rd.client_id as detail_client_id, rd.user_id, 
			rd.company_name_ar, rd.company_name_en, rd.commerical_number, 
			rd.activity_type, rd.address, rd.phone_number, rd.email, 
			r.qr_identifier,
			rd.identity_number, rd.mobile_number, rd.item_cost, rd.quantity,
			rd.for_official_use, rd.country_producer, rd.standard_of_origin,
			r.exporter_name,
			rd.extra,
			rd.foreign_items_cost
		FROM requests r
		JOIN request_details rd ON r.id = rd.request_id
		WHERE r.id = $1`

	var description models.NullString
	var status models.NullString
	var clientID, createdBy, requestTypeID models.NullInt64
	var createdAt, updatedAt time.Time
	fmt.Println("request-id:", id)
	err := r.db.QueryRow(query, id).Scan(
		&requestDetails.RequestID,
		&requestDetails.SerialNumber,
		&requestDetails.Title,
		&description,
		&status,
		&clientID,
		&createdBy,
		&requestTypeID,
		&createdAt,
		&updatedAt,
		&requestDetails.ID,
		&requestDetails.ClientName,
		&requestDetails.TransferDetail,
		&requestDetails.Signs,
		&requestDetails.NumberOfParcel,
		&requestDetails.Description,
		&requestDetails.Weight,
		&requestDetails.NetWeight,
		&requestDetails.InvoiceNumber,
		&requestDetails.InvoiceDate,
		&requestDetails.ClientID,
		&requestDetails.UserID,
		&requestDetails.CompanyName,
		&requestDetails.CompanyNameEn,
		&requestDetails.CommercialNumber,
		&requestDetails.ActivityType,
		&requestDetails.Address,
		&requestDetails.PhoneNumber,
		&requestDetails.Email,
		&requestDetails.QRIdentifier,
		&requestDetails.IdentityNumber,
		&requestDetails.MobileNumber,
		&requestDetails.Value,
		&requestDetails.Quantity,
		&requestDetails.ForOfficialUse,
		&requestDetails.CountryProducer,
		&requestDetails.StandardOfOrigin,
		&requestDetails.ExporterName,
		&requestDetails.Extra,
		&requestDetails.ForignItemCost,
	)

	if err != nil {
		return nil, err
	}

	return requestDetails, nil
}

// Update updates a request
func (r *RequestRepository) Update(request *models.Request) error {
	query := `
		UPDATE requests 
		SET title = $2, description = $3, status = $4, approved_by = $5, 
		    rejection_reason = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	var approvedBy sql.NullInt64
	if request.ApprovedByUserID != nil {
		approvedBy = sql.NullInt64{Int64: *request.ApprovedByUserID, Valid: true}
	}
	var rejection sql.NullString
	if request.RejectionReason != nil {
		rejection = sql.NullString{String: *request.RejectionReason, Valid: true}
	}

	err := r.db.QueryRow(
		query,
		request.ID,
		request.Title,
		request.Description,
		request.Status,
		approvedBy,
		rejection,
	).Scan(&request.UpdatedAt)

	return err
}

// Update request identifier
func (r *RequestRepository) UpdateIdentifier(request *models.Request) error {
	query := `
		UPDATE requests 
		SET qr_identifier = $2
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		request.ID,
		request.QRIdentifier,
	).Scan(&request.UpdatedAt)

	return err
}

// Approve approves a request
func (r *RequestRepository) Approve(requestID, approvedByStaffID int64) error {
	query := `
		UPDATE requests 
		SET status = 'APPROVED', approved_by = $2, approved_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = 'ISSUED'`

	result, err := r.db.Exec(query, requestID, approvedByStaffID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("request not found or not in pending status")
	}

	return nil
}

// Reject rejects a request
func (r *RequestRepository) Reject(requestID, approvedByStaffID int64, rejectionReason string) error {
	query := `
		UPDATE requests 
		SET status = 'REJECTED', approved_by = $2, rejection_reason = $3, 
		    rejected_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = 'ISSUED'`

	result, err := r.db.Exec(query, requestID, approvedByStaffID, rejectionReason)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("request not found or not in pending status")
	}

	return nil
}

func (r *RequestRepository) MarkAsPaid(requestID, approvedByStaffID int64, ratio float64) error {
	query := `
		UPDATE requests 
		SET status = 'PAID', approved_by = $2, approved_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP, ratio = $3
		WHERE id = $1 AND status = 'APPROVED'`

	result, err := r.db.Exec(query, requestID, approvedByStaffID, ratio)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("request not found or not in pending status")
	}

	return nil
}

// SumAllRequestRatios calculates the sum of ratios for all requests
func (r *RequestRepository) SumAllRequestRatios() (float64, error) {
	var sum float64
	query := `
		SELECT COALESCE(SUM(r.ratio), 0) 
		FROM requests r 
		WHERE r.status = 'PAID'`
	
	err := r.db.QueryRow(query).Scan(&sum)
	return sum, err
}

// SumRequestRatiosByBranch calculates the sum of ratios for requests in a specific branch
func (r *RequestRepository) SumRequestRatiosByBranch(branchID int64) (float64, error) {
	var sum float64
	query := `
		SELECT COALESCE(SUM(r.ratio), 0) 
		FROM requests r 
		WHERE r.branch_id = $1 AND r.status = 'PAID'`
	
	err := r.db.QueryRow(query, branchID).Scan(&sum)
	return sum, err
}

// SumRatiosPerBranch calculates the sum of ratios grouped by branch name
func (r *RequestRepository) SumRatiosPerBranch() ([]map[string]interface{}, error) {
	results := [] map[string]interface{}{}
	query := `
		SELECT b.name, COALESCE(SUM(r.ratio), 0) as ratio_sum
		FROM branches b
		LEFT JOIN requests r ON r.branch_id = b.id AND r.status = 'PAID'
		GROUP BY b.id, b.name`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branchName string
		var ratioSum float64
		if err := rows.Scan(&branchName, &ratioSum); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"branch_name": branchName,
			"ratio_sum":   ratioSum,
		})
	}

	return results, nil
}

func (r *RequestRepository) GenerateSerialNumber(requestID int64, branchID int64) error {
	query := `
		UPDATE requests 
		SET serial_number = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = 'PAID'`

	query2 := `
		SELECT code FROM branches WHERE id = $1`
	var code string
	err := r.db.QueryRow(query2, branchID).Scan(&code)
	if err != nil {
		return err
	}
	code += "-" + strconv.FormatInt(requestID, 10)
	result, err := r.db.Exec(query, requestID, code)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("request not found or not in pending status")
	}

	return nil
}

// CreateRequestDetail creates a new request detail record
func (r *RequestRepository) CreateRequestDetail(requestDetail *models.RequestDetail, requestID int64) error {
	query := `
		INSERT INTO request_details (
			client_name, transfer_detail, signs, number_of_parcel, request_id, 
			description, weight, net_weight, invoice_number, invoice_date, 
			client_id, user_id, company_name_ar, commerical_number, 
			activity_type, address, phone_number, email, identity_number, mobile_number,
			quantity, item_cost, for_official_use, country_producer, standard_of_origin,
			extra, foreign_items_cost
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25,
			$26,$27
		)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		requestDetail.ClientName,
		requestDetail.TransferDetail,
		requestDetail.Signs,
		requestDetail.NumberOfParcel,
		requestID,
		requestDetail.Description,
		requestDetail.Weight,
		requestDetail.NetWeight,
		requestDetail.InvoiceNumber,
		requestDetail.InvoiceDate,
		requestDetail.ClientID,
		requestDetail.UserID,
		requestDetail.CompanyName,
		requestDetail.CommercialNumber,
		requestDetail.ActivityType,
		requestDetail.Address,
		requestDetail.PhoneNumber,
		requestDetail.Email,
		requestDetail.IdentityNumber,
		requestDetail.MobileNumber,
		requestDetail.Quantity,
		requestDetail.Value,
		requestDetail.ForOfficialUse,
		requestDetail.CountryProducer,
		requestDetail.StandardOfOrigin,
		requestDetail.Extra,
		requestDetail.ForignItemCost,
	).Scan(&requestDetail.ID, &requestDetail.CreatedAt, &requestDetail.UpdatedAt)

	return err
}

// Delete deletes a request
func (r *RequestRepository) Delete(id int64) error {
	query := `DELETE FROM requests WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves requests with pagination and filtering active requests
func (r *RequestRepository) List(params models.FilterParams, pagination models.PaginationParams, userRole string) ([]*models.RequestResponse, int, error) {
	var requestsMap = make(map[int64]*models.RequestResponse)
	var orderedRequests []*models.RequestResponse
	var total int

	// Build WHERE clause for requests
	whereClause, args := r.buildWhereClause(params, userRole)
	fmt.Println("whereClause", whereClause)
	fmt.Println("args", args)
	// Count total distinct requests
	countQuery := "SELECT COUNT(DISTINCT r.id) FROM requests r JOIN request_details rd ON r.id = rd.request_id" + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderClause := r.buildOrderClause(params)

	// Main query with pagination to get requests with their details
	query := fmt.Sprintf(`
		SELECT 
			r.id as request_id, r.serial_number, r.title, r.description, r.status, r.client_id,
			r.created_by,r.request_type_id, r.created_at, r.updated_at,
			rd.id as detail_id, rd.client_name, rd.transfer_detail, rd.signs, 
			rd.number_of_parcel, rd.description as detail_description, rd.weight, rd.net_weight, 
			rd.invoice_number, rd.invoice_date, rd.client_id as detail_client_id, rd.user_id, 
			rd.company_name_ar, rd.company_name_en, rd.commerical_number, 
			rd.activity_type, rd.address, rd.phone_number, rd.email, 
			rd.identity_number, rd.mobile_number, rd.item_cost, rd.quantity,
			rt.name as request_type_name,
			rd.for_official_use, rd.country_producer, rd.standard_of_origin,r.exporter_name,r.ratio,
			rd.extra,
			rd.foreign_items_cost,
			r.branch_id,r.rejection_reason
		FROM requests r
		JOIN request_details rd ON r.id = rd.request_id
		JOIN request_types rt ON r.request_type_id = rt.id
		%s%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, len(args)+1, len(args)+2)
	fmt.Println(query)
	args = append(args, pagination.GetLimit(), pagination.GetOffset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var requestID int64
		var branchID int64
		var serialNumber, title, description models.NullString
		var status models.RequestStatus
		var clientID, detailClientID int64
		var createdBy models.NullInt64
		var requestTypeID *int64
		var createdAt, updatedAt time.Time
		var ratio *float64
		var rejectionReason models.NullString
		detail := &models.RequestDetailDto{}

		err := rows.Scan(
			&requestID,
			&serialNumber,
			&title,
			&description,
			&status,
			&clientID,
			&createdBy,
			&requestTypeID,
			&createdAt,
			&updatedAt,
			&detail.ID,
			&detail.ClientName,
			&detail.TransferDetail,
			&detail.Signs,
			&detail.NumberOfParcel,
			&detail.Description,
			&detail.Weight,
			&detail.NetWeight,
			&detail.InvoiceNumber,
			&detail.InvoiceDate,
			&detailClientID,
			&detail.UserID,
			&detail.CompanyName,
			&detail.CompanyNameEn,
			&detail.CommercialNumber,
			&detail.ActivityType,
			&detail.Address,
			&detail.PhoneNumber,
			&detail.Email,
			&detail.IdentityNumber,
			&detail.MobileNumber,
			&detail.Value,
			&detail.Quantity,
			&detail.RequestTypeName,
			&detail.ForOfficialUse,
			&detail.CountryProducer,
			&detail.StandardOfOrigin,
			&detail.ExporterName,
			&ratio,
			&detail.Extra,
			&detail.ForignItemCost,
			&branchID,
			&rejectionReason,
		)
		if err != nil {
			fmt.Println("Scan error:", err)
			return nil, 0, err
		}
		detail.SerialNumber = serialNumber
		// If we haven't seen this request before, create a new entry
		if _, exists := requestsMap[requestID]; !exists {
			var serialNumberPtr *string
			if serialNumber.Valid {
				str := serialNumber.String
				serialNumberPtr = &str
			}

			requestsMap[requestID] = &models.RequestResponse{
				ID:              requestID,
				SerialNumber:    serialNumberPtr,
				Title:           title.String,
				Description:     description.String,
				Status:          status,
				RequestTypeID:   requestTypeID,
				RejectionReason: &rejectionReason.String,
				RequestTypeName: &detail.RequestTypeName.String,
				CreatedAt:       createdAt,
				UpdatedAt:       updatedAt,
				BranchID:        &branchID,
				Ratio:           ratio,
				RequestDetails:  []*models.RequestDetailDto{},
			}
			orderedRequests = append(orderedRequests, requestsMap[requestID])
		}

		// Add the detail to the request's details
		detail.RequestID = requestID
		requestsMap[requestID].RequestDetails = append(requestsMap[requestID].RequestDetails, detail)
	}

	return orderedRequests, total, nil
}

// ListByClientID retrieves requests by client ID with pagination and filtering
func (r *RequestRepository) ListByClientID(clientID int64, params models.FilterParams, pagination models.PaginationParams, userRole string) ([]*models.RequestResponse, int, error) {
	// Delegate to ListDetails method which already provides grouped requests for a client
	return r.ListDetails(clientID, params, pagination, userRole)
}

func (r *RequestRepository) ListDetails(clientID int64, params models.FilterParams, pagination models.PaginationParams, userRole string) ([]*models.RequestResponse, int, error) {
	params.Search = "" // Override search to focus on client-specific filtering

	var requestsMap = make(map[int64]*models.RequestResponse)
	var orderedRequests []*models.RequestResponse
	var total int

	// Build WHERE clause with client filter
	whereClause, args := r.buildWhereClauseWithClient(params, clientID, userRole)

	// Count total distinct requests
	countQuery := "SELECT COUNT(DISTINCT r.id) FROM requests r JOIN request_details rd ON r.id = rd.request_id" + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderClause := r.buildOrderClause(params)

	// Single query to fetch both requests and their details
	query := fmt.Sprintf(`
        SELECT 
            r.id as request_id, r.created_by, r.created_at, r.updated_at, r.status,
            rd.id as detail_id, rd.client_name, rd.transfer_detail, rd.signs, 
            rd.number_of_parcel, rd.description, rd.weight, rd.net_weight, 
            rd.invoice_number, rd.invoice_date, rd.client_id, rd.user_id, 
            rd.company_name_ar, rd.company_name_en, rd.commerical_number, 
            rd.activity_type, rd.address, rd.phone_number, rd.email, 
            rd.identity_number, rd.mobile_number,
            rt.name as ServiceType
        FROM requests r
        JOIN request_details rd ON r.id = rd.request_id
        %s%s
        ORDER BY r.id %s, rd.id
        LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, orderClause, len(args)+1, len(args)+2)

	args = append(args, pagination.GetLimit(), pagination.GetOffset())

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var requestID int64
		var request models.Request
		var createdBy models.NullInt64
		detail := &models.RequestDetailDto{}

		err := rows.Scan(
			&requestID,
			&createdBy,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.Status,
			&detail.ID,
			&detail.ClientName,
			&detail.TransferDetail,
			&detail.Signs,
			&detail.NumberOfParcel,
			&detail.Description,
			&detail.Weight,
			&detail.NetWeight,
			&detail.InvoiceNumber,
			&detail.InvoiceDate,
			&detail.ClientID,
			&detail.UserID,
			&detail.CompanyName,
			&detail.CompanyNameEn,
			&detail.CommercialNumber,
			&detail.ActivityType,
			&detail.Address,
			&detail.PhoneNumber,
			&detail.Email,
			&detail.IdentityNumber,
			&detail.MobileNumber,
			&detail.Quantity,
			&detail.Value,
			&detail.ForOfficialUse,
			&detail.CountryProducer,
		)
		if err != nil {
			return nil, 0, err
		}

		// If we haven't seen this request before, create a new entry
		if _, exists := requestsMap[requestID]; !exists {
			requestsMap[requestID] = &models.RequestResponse{
				ID:             requestID,
				CreatedAt:      request.CreatedAt,
				UpdatedAt:      request.UpdatedAt,
				Status:         request.Status,
				RequestDetails: []*models.RequestDetailDto{},
			}
			orderedRequests = append(orderedRequests, requestsMap[requestID])
		}

		// Add the detail to the request's details
		detail.RequestID = requestID
		requestsMap[requestID].RequestDetails = append(requestsMap[requestID].RequestDetails, detail)
	}

	return orderedRequests, total, nil
}

// GetHistory retrieves request history
func (r *RequestRepository) GetHistory(requestID int64) ([]*models.RequestHistory, error) {
	var history []*models.RequestHistory

	query := `
		SELECT id, request_id, old_status, new_status, changed_by_user_id, change_reason, created_at
		FROM request_history 
		WHERE request_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		entry := &models.RequestHistory{}
		var oldStatus sql.NullString
		var changeReason sql.NullString
		err := rows.Scan(
			&entry.ID,
			&entry.RequestID,
			&oldStatus,
			&entry.NewStatus,
			&entry.ChangedByUserID,
			&changeReason,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if oldStatus.Valid {
			s := models.RequestStatus(oldStatus.String)
			entry.OldStatus = &s
		}
		if changeReason.Valid {
			s := changeReason.String
			entry.ChangeReason = &s
		}
		history = append(history, entry)
	}

	return history, nil
}

// buildWhereClause builds the WHERE clause for filtering
func (r *RequestRepository) buildWhereClause(params models.FilterParams, userRole string) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	fmt.Println(params)
	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d OR serial_number ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	// Apply role-based status filtering
	switch userRole {
	case models.RoleBranchAdmin:
		// Branch admins can view all statuses, so we don't add a status filter unless specified in params.
	case models.RoleAccountant:
		// Accountants can only view approved requests
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(models.StatusApproved))
		argIndex++
	case models.RoleAuditor:
		// Auditors can only view submitted/issued requests
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(models.StatusSubmitted))
		argIndex++
	default:
		// For other roles, apply any status filter from params
		if params.Status != "" {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, params.Status)
			argIndex++
		}
	}

	if params.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
			args = append(args, startDate)
			argIndex++
		}
	}

	if params.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
			args = append(args, endDate.Add(24*time.Hour))
			argIndex++
		}
	}

	// Apply branch filtering based on user role
	if userRole == models.RoleStaff || userRole == models.RoleBranchAdmin || userRole == models.RoleAccountant || userRole == models.RoleAuditor {
		// Staff, branch admin, and accountant users always see only their branch's requests
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argIndex))
		args = append(args, params.BranchID)
		argIndex++
	} else if userRole == models.RoleAdmin && params.BranchIDFilter != nil && *params.BranchIDFilter != 0 {
		// Admin can filter by specific branch if requested
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argIndex))
		args = append(args, params.BranchIDFilter)
		argIndex++
	}

	if params.Serial != "" {
		conditions = append(conditions, fmt.Sprintf("serial_number = $%d", argIndex))
		args = append(args, params.Serial)
		argIndex++
	}

	var whereClause string
	var whereConditions []string

	// Add conditions for branch filtering (only if not admin)
	whereConditions = append(whereConditions, conditions...)

	// Construct final WHERE clause
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	return whereClause, args
}

// buildWhereClauseWithClient builds the WHERE clause with client filter
func (r *RequestRepository) buildWhereClauseWithClient(params models.FilterParams, clientID int64, userRole string) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Always filter by client ID
	conditions = append(conditions, fmt.Sprintf("r.client_id = $%d", argIndex))
	args = append(args, clientID)
	argIndex++

	// Add branch filter based on user role
	if userRole == models.RoleStaff || userRole == models.RoleBranchAdmin || userRole == models.RoleAccountant || userRole == models.RoleAuditor {
		// Staff, branch admin, and accountant users always see only their branch's requests
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argIndex))
		args = append(args, params.BranchID)
		argIndex++
	} else if userRole != models.RoleAdmin && params.BranchID != 0 {
		// Other non-admin users with branch filters
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argIndex))
		args = append(args, params.BranchID)
		argIndex++
	}

	// Apply role-based status filtering
	switch userRole {
	case models.RoleAccountant:
		// Accountants can only view approved requests
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(models.StatusApproved))
		argIndex++
	case models.RoleAuditor:
		// Auditors can only view submitted/issued requests
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(models.StatusSubmitted))
		argIndex++
	default:
		// For other roles, apply any status filter from params
		if params.Status != "" {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, params.Status)
			argIndex++
		}
	}

	if params.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
			args = append(args, startDate)
			argIndex++
		}
	}

	if params.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
			args = append(args, endDate.Add(24*time.Hour))
			argIndex++
		}
	}

	whereClause := " WHERE " + strings.Join(conditions, " AND ")
	return whereClause, args
}

// buildOrderClause builds the ORDER BY clause
func (r *RequestRepository) buildOrderClause(params models.FilterParams) string {
	sortBy := "created_at"
	if params.SortBy != "" {
		switch params.SortBy {
		case "title", "status", "created_at", "updated_at", "serial_number":
			sortBy = params.SortBy
		}
	}

	sortOrder := "DESC"
	if params.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}

func (r *RequestRepository) ListDetailsCount(clientID int64, params models.FilterParams, userRole string) (int, error) {
	params.Search = "" // Override search to focus on client-specific filtering

	var total int

	// Build WHERE clause with client filter
	whereClause, args := r.buildWhereClause(params, userRole)

	// Count total distinct requests
	countQuery := "SELECT COUNT(DISTINCT r.id)  FROM requests r JOIN request_details rd ON r.id = rd.request_id" + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetBySerial retrieves a request by its serial number
func (r *RequestRepository) GetBySerial(serialNumber string) (*models.Request, error) {
	query := `
		SELECT id, serial_number, title, description, status, client_id, created_by, 
		       approved_by, rejection_reason, approved_at, rejected_at, request_type_id, 
		       branch_id, exporter_name, qr_identifier, created_at, updated_at
		FROM requests 
		WHERE serial_number = $1
	`

	var request models.Request
	var approvedAt, rejectedAt sql.NullTime

	err := r.db.QueryRow(query, serialNumber).Scan(
		&request.ID,
		&request.SerialNumber,
		&request.Title,
		&request.Description,
		&request.Status,
		&request.ClientID,
		&request.CreatedByUserID,
		&request.ApprovedByUserID,
		&request.RejectionReason,
		&approvedAt,
		&rejectedAt,
		&request.RequestTypeID,
		&request.BranchID,
		&request.ExporterName,
		&request.QRIdentifier,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Handle nullable time fields
	if approvedAt.Valid {
		request.ApprovedAt = &approvedAt.Time
	}
	if rejectedAt.Valid {
		request.RejectedAt = &rejectedAt.Time
	}

	return &request, nil
}
