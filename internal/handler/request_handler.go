package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"Chumber-Workflow-System/internal/middleware"
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/repository"
	"Chumber-Workflow-System/internal/service"

	"Chumber-Workflow-System/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RequestHandler handles request-related HTTP requests
type RequestHandler struct {
	requestService *service.RequestService
	activityRepo   *repository.UserActivityRepository
	validator      *validator.Validate
	config         *config.Config
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(requestService *service.RequestService, activityRepo *repository.UserActivityRepository, cfg config.Config) *RequestHandler {
	return &RequestHandler{
		requestService: requestService,
		activityRepo:   activityRepo,
		validator:      validator.New(),
		config:         &cfg,
	}
}

// CreateRequest handles creating a new request
// @Summary Create a new request
// @Description Create a new request (clients create for themselves, staff can create for clients)
// @Tags requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateRequestRequest true "Create request"
// @Success 201 {object} models.APIResponse{data=models.RequestResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /requests [post]
func (h *RequestHandler) CreateRequest(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	username, exists := middleware.GetCurrentUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	// userID is already an int64

	userRole, exists := middleware.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	var req models.CreateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}
	if userRole == models.RoleClient {
		req.ClientID = userID
	}
	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	request, err := h.requestService.CreateRequest(&req, userID, userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	// Track user activity in a background goroutine
	h.trackUserActivity(userID, username, "create", "request", request.ID, fmt.Sprintf("قام الموظف %s بانشاء طلب جديد رقم %d", username, request.ID))

	c.JSON(http.StatusCreated, models.SuccessResponse("Request created successfully", request.ToResponse()))
}

func (h *RequestHandler) AcceptRequest(c *gin.Context) {
	currentUserID, exists := middleware.GetCurrentUserID(c)
	userRole, exists := middleware.GetCurrentUserRole(c)
	username, exists := middleware.GetCurrentUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	idStr := c.Param("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	request, err := h.requestService.ApproveRequest(requestID, currentUserID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}
	h.trackUserActivity(currentUserID, "", "accept", "request", requestID, fmt.Sprintf("%sقام الموظف %d بالقبول عن الطلب رقم %d", username, currentUserID, requestID))
	c.JSON(http.StatusOK, models.SuccessResponse("Request accepted successfully", request.ToResponse()))
}

// GetRequest handles getting a request by ID
// @Summary Get request by ID
// @Description Get a request by ID (clients can only see their own, staff/admin can see all)
// @Tags requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} models.APIResponse{data=models.RequestResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /requests/{id} [get]
func (h *RequestHandler) GetRequest(c *gin.Context) {
	currentUserID, exists := middleware.GetCurrentUserID(c)
	userRole, exists := middleware.GetCurrentUserRole(c)
	username, exists := middleware.GetCurrentUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	idStr := c.Param("id")
	fmt.Println("idStr", idStr)
	requestID, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	requestDetails, err := h.requestService.GetRequestWithDetails(requestID, currentUserID, userRole, h.config.JWT.Secret)

	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			return
		}
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}
	h.trackUserActivity(currentUserID, "", "get", "request", requestID, fmt.Sprintf("قام الموظف %s بالاستعلام عن الطلب رقم %d", username, requestID))
	c.JSON(http.StatusOK, models.SuccessResponse("Request retrieved", requestDetails))
}

// GetRequestByIdentifier handles getting a request by encrypted identifier (unauthenticated)
// @Summary Get request by encrypted identifier
// @Description Get a request by encrypted identifier (only approved requests)
// @Tags requests
// @Produce json
// @Param identifier path string true "Encrypted Request Identifier"
// @Success 200 {object} models.APIResponse{data=models.RequestResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /request/{identifier} [get]
func (h *RequestHandler) GetRequestByIdentifier(c *gin.Context) {
	encryptedID := c.Param("identifier")
	if encryptedID == "" {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid identifier"))
		return
	}

	// Use the new service method to get request with details using identifier
	requestDetails, err := h.requestService.GetRequestWithDetailsUsingIdentifier(encryptedID, h.config.JWT.Secret)
	if err != nil {
		if err.Error() == "request not found" || err.Error() == "request not approved" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid identifier"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Request retrieved", requestDetails))
}

// GetRequestBySerial handles getting a request by serial number (unauthenticated)
// @Summary Get request by serial number
// @Description Get a request by serial number (only approved requests)
// @Tags requests
// @Produce json
// @Param serial path string true "Request Serial Number"
// @Success 200 {object} models.APIResponse{data=models.RequestResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /request/serial/{serial} [get]
func (h *RequestHandler) GetRequestBySerial(c *gin.Context) {
	serialNumber := c.Param("serial")
	if serialNumber == "" {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid serial number"))
		return
	}

	// Get request by serial number
	request, err := h.requestService.GetRequestBySerial(serialNumber)
	if err != nil {
		if err.Error() == "request not found" || err.Error() == "request not approved" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Request retrieved", request.ToResponse()))
}

// UpdateRequest handles updating a request
// @Summary Update a request
// @Description Update a request (only pending requests can be updated)
// @Tags requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Param request body models.UpdateRequestRequest true "Update request"
// @Success 200 {object} models.APIResponse{data=models.RequestResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /requests/{id} [put]
func (h *RequestHandler) UpdateRequest(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	userID := userIDStr

	userRole, exists := middleware.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	idStr := c.Param("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	var req models.UpdateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body"))
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	request, err := h.requestService.UpdateRequest(requestID, &req, userID, userRole)
	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			return
		}
		fmt.Println("err", err)
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Request updated successfully", request.ToResponse()))
}

// ListRequests handles listing requests with pagination
// @Summary List requests
// @Description List requests with pagination and filtering (clients see only their own, staff/admin see all)
// @Tags requests
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Param status query string false "Filter by status"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Router /requests [get]
func (h *RequestHandler) ListRequests(c *gin.Context) {
	fmt.Println("ListRequests")
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}
	fmt.Println("userIDStr", userIDStr)
	userID := userIDStr

	userRole, exists := middleware.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}
	fmt.Println("userRole", userRole)
	var branchID int64
	if userRole != models.RoleAdmin {
		branchID, exists = middleware.GetCurrentUserBranchID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
			return
		}
	}

	fmt.Println("branchID", branchID)

	// Parse pagination parameters
	pagination := models.GetDefaultPagination()
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}
	fmt.Println("pagination", pagination)
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}
	fmt.Println("pagination", pagination)
	branchIDfilter := int64(0)
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" && userRole == models.RoleAdmin {
		if branchIDInt, err := strconv.Atoi(branchIDStr); err == nil && branchIDInt > 0 {
			branchIDfilter = int64(branchIDInt)
		}
	}

	// Parse filter parameters
	params := models.FilterParams{
		Search:         c.Query("search"),
		Status:         c.Query("status"),
		StartDate:      c.Query("start_date"),
		EndDate:        c.Query("end_date"),
		SortBy:         c.Query("sort_by"),
		SortOrder:      c.Query("sort_order"),
		Serial:         c.Query("serial_number"),
		BranchID:       int64(branchID),
		BranchIDFilter: &branchIDfilter,
	}

	requests, paginationResp, err := h.requestService.ListRequests(params, pagination, userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	response := map[string]interface{}{
		"requests":   requests,
		"pagination": paginationResp,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Requests retrieved", response))
}

// ListRequests handles listing requests with pagination
// @Summary List requests
// @Description List requests with pagination and filtering (clients see only their own, staff/admin see all)
// @Tags requests
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Param status query string false "Filter by status"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Router /requests/count [get]
func (h *RequestHandler) ListRequestsCount(c *gin.Context) {
	fmt.Println("ListRequests")
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}
	fmt.Println("userIDStr", userIDStr)
	userID := userIDStr

	userRole, exists := middleware.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}
	fmt.Println("userRole", userRole)
	branchID, exists := middleware.GetCurrentUserBranchID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}
	fmt.Println("branchID", branchID)

	// Parse filter parameters
	params := models.FilterParams{
		Search:    c.Query("search"),
		Status:    c.Query("status"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		BranchID:  int64(branchID),
	}

	response, err := h.requestService.ListRequestsCount(params, userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Requests retrieved", response))
}

// ApproveRequest handles approving a request (staff/admin only)
// @Summary Approve a request
// @Description Approve a request and generate serial number (staff/admin only)
// @Tags requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Param request body models.ApproveRequestRequest true "Approve request"
// @Success 200 {object} models.APIResponse{data=models.RequestResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /requests/{id}/approve [post]
func (h *RequestHandler) ApproveRequest(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	userRole, exists := middleware.GetCurrentUserRole(c)
	username, exists := middleware.GetCurrentUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	idStr := c.Param("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	request, err := h.requestService.ApproveRequest(requestID, userID, userRole)
	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		if err.Error() == "only staff and admin can approve requests" {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	// Track user activity in a background goroutine
	h.trackUserActivity(userID, username, "approve", "request", requestID, fmt.Sprintf("قام الموظف %s بالموافقة على الطلب رقم %d", username, requestID))

	c.JSON(http.StatusOK, models.SuccessResponse("Request approved successfully", request.ToResponse()))
}

// RejectRequest handles rejecting a request (staff/admin only)
// @Summary Reject a request
// @Description Reject a request with reason (staff/admin only)
// @Tags requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Param request body models.RejectRequestRequest true "Reject request"
// RejectRequestBody defines the structure for rejection request payload
type RejectRequestBody struct {
	Reason string `json:"reason"`
}

// @Success 200 {object} models.APIResponse{data=models.RequestResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /requests/{id}/reject [post]

// trackUserActivity asynchronously creates a user activity record
func (h *RequestHandler) trackUserActivity(userID int64, username string, action string, module string, entityID int64, description string) {
	go func() {
		entityIDPtr := new(int64)
		*entityIDPtr = entityID

		activity := &models.UserActivities{
			UserId:      userID,
			Username:    username,
			Action:      action,
			Module:      module,
			EntityId:    entityIDPtr,
			Description: description,
		}

		err := h.activityRepo.Create(activity)
		if err != nil {
			fmt.Printf("Error recording user activity: %v\n", err)
		}
	}()
}

func (h *RequestHandler) RejectRequest(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	userRole, exists := middleware.GetCurrentUserRole(c)
	username, exists := middleware.GetCurrentUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	// Parse request body to get rejection reason
	var requestBody RejectRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	idStr := c.Param("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	// Pass the rejection reason to the service
	request, err := h.requestService.RejectRequest(requestID, userID, userRole, requestBody.Reason)
	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		if err.Error() == "only staff and admin can reject requests or branch admin" {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	// Track user activity in a background goroutine
	h.trackUserActivity(userID, username, "reject", "request", requestID, fmt.Sprintf("قام الموظف %s برفض الطلب رقم %d", username, requestID))

	c.JSON(http.StatusOK, models.SuccessResponse("Request rejected successfully", request.ToResponse()))
}
func (h *RequestHandler) MarkAsPaid(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	userRole, exists := middleware.GetCurrentUserRole(c)
	username, exists := middleware.GetCurrentUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	idStr := c.Param("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	request, err := h.requestService.MarkAsPaid(requestID, userID, userRole, h.config.JWT.Secret)
	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		if err.Error() == "only staff and admin can mark requests as paid" {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	// Track user activity in a background goroutine
	h.trackUserActivity(userID, username, "mark_paid", "request", requestID, fmt.Sprintf("قام الموظف %s بتعليم الطلب كمدفوع رقم %d", username, requestID))

	c.JSON(http.StatusOK, models.SuccessResponse("Request marked as paid successfully", request.ToResponse()))
}

// GetRatioSum handles getting the sum of ratios based on user permissions
// @Summary Get ratio sum
// @Description Get sum of ratios - all request ratios for admin, branch request ratios for branch admin
// @Tags requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=map[string]float64}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /requests/ratio-sum [get]
func (h *RequestHandler) GetRatioSum(c *gin.Context) {
	var branchID *int64
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if id, err := strconv.ParseInt(branchIDStr, 10, 64); err == nil {
			branchID = &id
		} else {
			c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid branch_id parameter"))
			return
		}
	}

	// Get ratio sum based on user permissions
	// Get ratio sum with optional branch filtering
	sum, err := h.requestService.GetRatioSumWithBranchFilter(branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}
	response := map[string]float64{
		"ratio_sum": sum,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Ratio sum retrieved successfully", response))
}

// GetRatioSumPerBranch handles getting the sum of ratios grouped by branch name (admin only)
// @Summary Get ratio sum per branch
// @Description Get sum of ratios grouped by branch name (admin only)
// @Tags requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=map[string]float64}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /requests/ratio-sum-per-branch [get]
func (h *RequestHandler) GetRatioSumPerBranch(c *gin.Context) {

	// Get ratio sum per branch
	ratioSumPerBranch, err := h.requestService.GetRatioSumPerBranch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Ratio sum per branch retrieved successfully", ratioSumPerBranch))
}

// GetRequestHistory handles getting request history
// @Summary Get request history
// @Description Get the history of status changes for a request
// @Tags requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} models.APIResponse{data=[]models.RequestHistory}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /requests/{id}/history [get]
func (h *RequestHandler) GetRequestHistory(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	userID := userIDStr

	userRole, exists := middleware.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	idStr := c.Param("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	history, err := h.requestService.GetRequestHistory(requestID, userID, userRole)
	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Request history retrieved", history))
}

// DeleteRequest handles deleting a request (admin only)
// @Summary Delete a request
// @Description Delete a request (admin only, only pending requests)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/requests/{id} [delete]
func (h *RequestHandler) DeleteRequest(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	userID := userIDStr

	userRole, exists := middleware.GetCurrentUserRole(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedErrorResponse())
		return
	}

	idStr := c.Param("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid request ID"))
		return
	}

	err = h.requestService.DeleteRequest(requestID, userID, userRole)
	if err != nil {
		if err.Error() == "request not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Request"))
			return
		}
		if err.Error() == "only admin can delete requests" {
			c.JSON(http.StatusForbidden, models.ForbiddenErrorResponse())
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Request deleted successfully", nil))
}
