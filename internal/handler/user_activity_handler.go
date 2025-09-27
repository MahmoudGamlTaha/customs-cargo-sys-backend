package handler

import (
	"Chumber-Workflow-System/internal/middleware"
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserActivityHandler handles user activity-related HTTP requests
type UserActivityHandler struct {
	service *service.UserActivityService
}

// NewUserActivityHandler creates a new user activity handler
func NewUserActivityHandler(service *service.UserActivityService) *UserActivityHandler {
	return &UserActivityHandler{
		service: service,
	}
}

// GetAllActivities gets all user activities with pagination
// @Summary List all user activities
// @Description Get all user activities with pagination
// @Tags activities
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.APIResponse{data=[]models.UserActivitiesResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /activities [get]
func (h *UserActivityHandler) GetAllActivities(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	pagination := getPaginationParams(c)

	activities, total, err := h.service.GetAll(userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []*models.UserActivitiesResponse
	for _, activity := range activities {
		response = append(response, activity.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"total":     total,
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})
}

// GetActivitiesByUserID gets user activities by user ID
// @Summary Get activities by user ID
// @Description Get user activities filtered by user ID with pagination
// @Tags activities
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.APIResponse{data=[]models.UserActivitiesResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /activities/user/{user_id} [get]
func (h *UserActivityHandler) GetActivitiesByUserID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pagination := getPaginationParams(c)

	activities, total, err := h.service.GetByUserID(userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []*models.UserActivitiesResponse
	for _, activity := range activities {
		response = append(response, activity.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"total":     total,
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})
}

// GetActivitiesByModuleAndEntityID gets user activities by module and entity ID
// @Summary Get activities by module and entity ID
// @Description Get user activities filtered by module and entity ID with pagination
// @Tags activities
// @Security BearerAuth
// @Produce json
// @Param module path string true "Module name"
// @Param entity_id path int true "Entity ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.APIResponse{data=[]models.UserActivitiesResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /activities/module/{module}/entity/{entity_id} [get]
func (h *UserActivityHandler) GetActivitiesByModuleAndEntityID(c *gin.Context) {
	module := c.Param("module")
	entityID, err := strconv.ParseInt(c.Param("entity_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	pagination := getPaginationParams(c)

	activities, total, err := h.service.GetByModuleAndEntityID(module, entityID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []*models.UserActivitiesResponse
	for _, activity := range activities {
		response = append(response, activity.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"total":     total,
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})
}

// Helper function to get pagination parameters from the request
func getPaginationParams(c *gin.Context) models.PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}
