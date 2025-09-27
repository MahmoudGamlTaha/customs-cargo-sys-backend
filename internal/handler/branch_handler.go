package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"
)

// BranchHandler handles branch-related HTTP requests
type BranchHandler struct {
	branchService *service.BranchService
	validator     *validator.Validate
}

// NewBranchHandler creates a new branch handler
func NewBranchHandler(branchService *service.BranchService) *BranchHandler {
	return &BranchHandler{
		branchService: branchService,
		validator:     validator.New(),
	}
}

// CreateBranch handles creating a new branch (admin only)
// @Summary Create a new branch
// @Description Create a new branch (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateBranchRequest true "Create branch request"
// @Success 201 {object} models.APIResponse{data=models.BranchResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /admin/branches [post]
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req models.CreateBranchRequest
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

	branch, err := h.branchService.CreateBranch(&req)
	if err != nil {
		if err.Error() == "branch code already exists" || err.Error() == "branch name already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Branch created successfully", branch.ToResponse()))
}

// GetBranch handles getting a branch by ID
// @Summary Get branch by ID
// @Description Get a branch by ID
// @Tags branches
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Success 200 {object} models.APIResponse{data=models.BranchResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /branches/{id} [get]
func (h *BranchHandler) GetBranch(c *gin.Context) {
	idStr := c.Param("id")
	branchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid branch ID"))
		return
	}

	branch, err := h.branchService.GetBranch(branchID)
	if err != nil {
		if err.Error() == "branch not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Branch"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Branch retrieved", branch.ToResponse()))
}

// ListBranches handles listing branches with pagination
// @Summary List branches
// @Description List branches with pagination and filtering
// @Tags branches
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Router /branches [get]
func (h *BranchHandler) ListBranches(c *gin.Context) {
	// Parse pagination parameters
	pagination := models.GetDefaultPagination()
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}

	// Parse filter parameters
	params := models.FilterParams{
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	branches, paginationResp, err := h.branchService.ListBranches(params, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	// Convert to response format
	branchResponses := make([]*models.BranchResponse, len(branches))
	for i, branch := range branches {
		branchResponses[i] = branch.ToResponse()
	}

	response := map[string]interface{}{
		"branches":   branchResponses,
		"pagination": paginationResp,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Branches retrieved", response))
}

// GetAllBranches handles getting all branches (for dropdown lists)
// @Summary Get all branches
// @Description Get all branches for dropdown lists
// @Tags branches
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.BranchResponse}
// @Failure 401 {object} models.APIResponse
// @Router /branches/all [get]
func (h *BranchHandler) GetAllBranches(c *gin.Context) {
	branches, err := h.branchService.GetAllBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	// Convert to response format
	branchResponses := make([]*models.BranchResponse, len(branches))
	for i, branch := range branches {
		branchResponses[i] = branch.ToResponse()
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Branches retrieved", branchResponses))
}

// UpdateBranch handles updating a branch (admin only)
// @Summary Update a branch
// @Description Update a branch (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Param request body models.UpdateBranchRequest true "Update branch request"
// @Success 200 {object} models.APIResponse{data=models.BranchResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/branches/{id} [put]
func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	idStr := c.Param("id")
	branchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid branch ID"))
		return
	}

	var req models.UpdateBranchRequest
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

	branch, err := h.branchService.UpdateBranch(branchID, &req)
	if err != nil {
		if err.Error() == "branch not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Branch"))
			return
		}
		if err.Error() == "branch code already exists" || err.Error() == "branch name already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Branch updated successfully", branch.ToResponse()))
}

// DeleteBranch handles deleting a branch (admin only)
// @Summary Delete a branch
// @Description Delete a branch (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/branches/{id} [delete]
func (h *BranchHandler) DeleteBranch(c *gin.Context) {
	idStr := c.Param("id")
	branchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid branch ID"))
		return
	}

	err = h.branchService.DeleteBranch(branchID)
	if err != nil {
		if err.Error() == "branch not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Branch"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Branch deleted successfully", nil))
}

