package handler

import (
	"net/http"
	"strconv"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	
)

// CompanyHandler handles company-related HTTP requests
type CompanyHandler struct {
	companyService *service.CompanyService
	validator      *validator.Validate
}

// NewCompanyHandler creates a new company handler
func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		validator:      validator.New(),
	}
}

// CreateCompany handles creating a new company (admin only)
// @Summary Create a new company
// @Description Create a new company (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateCompanyRequest true "Create company request"
// @Success 201 {object} models.APIResponse{data=models.CompanyResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /admin/companies [post]
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req models.CreateCompanyRequest
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

	company, err := h.companyService.CreateCompany(&req)
	if err != nil {
		if err.Error() == "company code already exists" || err.Error() == "company name already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Company created successfully", company.ToResponse()))
}

// GetCompany handles getting a company by ID
// @Summary Get company by ID
// @Description Get a company by ID
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} models.APIResponse{data=models.CompanyResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	idStr := c.Param("id")
	companyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid company ID"))
		return
	}

	company, err := h.companyService.GetCompany(companyID)
	if err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Company"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Company retrieved", company.ToResponse()))
}

// ListCompanies handles listing companies with pagination
// @Summary List companies
// @Description List companies with pagination and filtering
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Router /companies [get]
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
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

	companies, paginationResp, err := h.companyService.ListCompanies(params, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	// Convert to response format
	companyResponses := make([]*models.CompanyResponse, len(companies))
	for i, company := range companies {
		companyResponses[i] = company.ToResponse()
	}

	response := map[string]interface{}{
		"companies":  companyResponses,
		"pagination": paginationResp,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Companies retrieved", response))
}

// GetAllCompanies handles getting all companies (for dropdown lists)
// @Summary Get all companies
// @Description Get all companies for dropdown lists
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.CompanyResponse}
// @Failure 401 {object} models.APIResponse
// @Router /companies/all [get]
func (h *CompanyHandler) GetAllCompanies(c *gin.Context) {
	companies, err := h.companyService.GetAllCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	// Convert to response format
	companyResponses := make([]*models.CompanyResponse, len(companies))
	for i, company := range companies {
		companyResponses[i] = company.ToResponse()
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Companies retrieved", companyResponses))
}

// UpdateCompany handles updating a company (admin only)
// @Summary Update a company
// @Description Update a company (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Param request body models.UpdateCompanyRequest true "Update company request"
// @Success 200 {object} models.APIResponse{data=models.CompanyResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	idStr := c.Param("id")
	companyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid company ID"))
		return
	}

	var req models.UpdateCompanyRequest
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

	company, err := h.companyService.UpdateCompany(companyID, &req)
	if err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Company"))
			return
		}
		if err.Error() == "company code already exists" || err.Error() == "company name already exists" {
			c.JSON(http.StatusConflict, models.ConflictErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Company updated successfully", company.ToResponse()))
}

// DeleteCompany handles deleting a company (admin only)
// @Summary Delete a company
// @Description Delete a company (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /admin/companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	idStr := c.Param("id")
	companyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid company ID"))
		return
	}

	err = h.companyService.DeleteCompany(companyID)
	if err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Company"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Company deleted successfully", nil))
}
