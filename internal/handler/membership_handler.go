package handler

import (
	"net/http"
	"strconv"

	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MembershipHandler handles membership-related HTTP requests
type MembershipHandler struct {
	membershipService *service.MembershipService
	validator         *validator.Validate
}

// NewMembershipHandler creates a new membership handler
func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{
		membershipService: membershipService,
		validator:         validator.New(),
	}
}

// GetMembershipByID handles getting a membership by ID (unauthenticated)
// @Summary Get membership by ID
// @Description Get a membership by ID (only approved memberships)
// @Tags memberships
// @Produce json
// @Param id path string true "Membership ID"
// @Success 200 {object} models.APIResponse{data=models.MembershipsDTO}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /membership/check/{id} [get]
func (h *MembershipHandler) GetMembershipByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid membership ID"))
		return
	}

	// Convert string ID to int64
	membershipID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestErrorResponse("Invalid membership ID format"))
		return
	}

	// Get membership by ID
	membership, err := h.membershipService.GetMembershipByID(membershipID)
	if err != nil {
		if err.Error() == "membership not found" || err.Error() == "membership not approved" {
			c.JSON(http.StatusNotFound, models.NotFoundErrorResponse("Membership"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Membership retrieved", membership))
}
