package service

import (
	"Chumber-Workflow-System/internal/models"
	"Chumber-Workflow-System/pkg/client"
)

// MembershipService handles membership business logic
type MembershipService struct {
	membershipClient *client.MembershipClient
}

// NewMembershipService creates a new membership service
func NewMembershipService() *MembershipService {
	return &MembershipService{
		membershipClient: client.NewMembershipClient(),
	}
}

// GetMembershipByID retrieves a membership by its ID from external API
func (s *MembershipService) GetMembershipByID(membershipID int64) (*models.Memberships, error) {
	// Fetch membership from external API
	membership, err := s.membershipClient.GetMembershipByID(membershipID)
	if err != nil {
		return nil, err
	}

	return membership, nil
}
