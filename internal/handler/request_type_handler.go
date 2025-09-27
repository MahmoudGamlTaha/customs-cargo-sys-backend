package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"

	"github.com/gin-gonic/gin"
)

// RequestTypeHandler handles request type-related HTTP requests
type RequestTypeHandler struct {
	requestTypeService *service.RequestTypeService
}

// NewRequestTypeHandler creates a new request type handler
func NewRequestTypeHandler(requestTypeService *service.RequestTypeService) *RequestTypeHandler {
	return &RequestTypeHandler{requestTypeService: requestTypeService}
}

// ListRequestTypes handles listing request types with pagination
// @Summary List request types
// @Description List request types with pagination and filtering
// @Tags request-types
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field (name, price, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Router /request-types [get]
func (h *RequestTypeHandler) ListRequestTypes(c *gin.Context) {
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

	items, paginationResp, err := h.requestTypeService.ListRequestTypes(params, pagination)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	// Convert to response format
	responses := make([]*models.RequestTypeResponse, len(items))
	for i, item := range items {
		responses[i] = item.ToResponse()
	}

	response := map[string]interface{}{
		"request_types": responses,
		"pagination":    paginationResp,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Request types retrieved", response))
}

// GetAllRequestTypes handles getting all request types (for dropdowns)
// @Summary Get all request types
// @Description Get all request types for dropdown lists
// @Tags request-types
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.RequestTypeResponse}
// @Failure 401 {object} models.APIResponse
// @Router /request-types/all [get]
func (h *RequestTypeHandler) GetAllRequestTypes(c *gin.Context) {
	items, err := h.requestTypeService.GetAllRequestTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	responses := make([]*models.RequestTypeResponse, len(items))
	for i, item := range items {
		responses[i] = item.ToResponse()
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Request types retrieved", responses))
}
