package client

import (
	"Chumber-Workflow-System/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// MembershipClient handles external API requests for membership data
type MembershipClient struct {
	baseURL    string
	httpClient *http.Client
}

// MembershipData represents the membership data structure from the API
type MembershipData struct {
	ID               int64      `json:"id"`
	MembershipCode   string     `json:"membership_code"`
	MemberName       string     `json:"member_name"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone"`
	Status           string     `json:"status"`
	MembershipType   string     `json:"membership_type"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedByUserID  int64      `json:"created_by_user_id"`
	ApprovedAt       *time.Time `json:"approved_at"`
	ApprovedByUserID *int64     `json:"approved_by_user_id"`
	RejectionReason  *string    `json:"rejection_reason"`
	RejectedAt       *time.Time `json:"rejected_at"`
}

// NewMembershipClient creates a new membership client instance
func NewMembershipClient() *MembershipClient {
	baseURL := os.Getenv("EXTERNAL_API_BASE_URL")
	if baseURL == "" {
		// Default fallback - should be set in environment
		baseURL = "http://localhost:8080"
	}

	return &MembershipClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetMembershipByID fetches a membership by ID from the external API
func (c *MembershipClient) GetMembershipByID(membershipID int64) (*models.Memberships, error) {
	// Construct the API endpoint URL
	url := fmt.Sprintf("%s/api/v1/membership/check/%d", c.baseURL, membershipID)

	// Make HTTP GET request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Handle non-success responses
	if !apiResponse.Success {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("membership not found")
		}

		errorMsg := "unknown error"
		if apiResponse.Error != nil {
			errorMsg = apiResponse.Error.Message
		} else if apiResponse.Message != "" {
			errorMsg = apiResponse.Message
		}

		return nil, fmt.Errorf("API error: %s", errorMsg)
	}

	// Check if data exists
	if apiResponse.Data == nil {
		return nil, fmt.Errorf("no data in API response")
	}

	// Convert the data to MembershipData struct
	dataBytes, err := json.Marshal(apiResponse.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal API data: %w", err)
	}

	var membershipData models.Memberships
	if err := json.Unmarshal(dataBytes, &membershipData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal membership data: %w", err)
	}

	return &membershipData, nil
}
