package client

import (
	"Chumber-Workflow-System/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// RequestClient handles external API requests for request data
type RequestClient struct {
	baseURL    string
	httpClient *http.Client
}

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// APIError represents error structure in API responses
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RequestData represents the request data structure from the API
type RequestData struct {
	ID                 int64     `json:"id"`
	SerialNumber       string    `json:"serial_number"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	RejectionReason    *string   `json:"rejection_reason"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ApprovedAt         *time.Time `json:"approved_at"`
	RejectedAt         *time.Time `json:"rejected_at"`
	CreateByUserID     int64     `json:"create_by_user_id"`
	ApprovedByUserID   *int64    `json:"approved_by_user_id"`
	RequestTypeID      int64     `json:"request_type_id"`
	RequestDetails     interface{} `json:"request_details"`
	QRIdentifier       string    `json:"qr_identifier"`
}

// NewRequestClient creates a new request client instance
func NewRequestClient() *RequestClient {
	baseURL := os.Getenv("EXTERNAL_API_BASE_URL")
	if baseURL == "" {
		// Default fallback - should be set in environment
		baseURL = "http://localhost:8080"
	}

	return &RequestClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetRequestBySerial fetches a request by serial number from the external API
func (c *RequestClient) GetRequestBySerial(serialNumber string) (*models.Request, error) {
	// Construct the API endpoint URL
	url := fmt.Sprintf("%s/api/v1/request/serial/%s", c.baseURL, serialNumber)

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
			return nil, fmt.Errorf("request not found")
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

	// Convert the data to RequestData struct
	dataBytes, err := json.Marshal(apiResponse.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal API data: %w", err)
	}

	var requestData RequestData
	if err := json.Unmarshal(dataBytes, &requestData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	// Map the API response to our internal Request model
	request := c.mapToRequestModel(&requestData)

	return request, nil
}

// mapToRequestModel converts RequestData from API to internal Request model
func (c *RequestClient) mapToRequestModel(data *RequestData) *models.Request {
	// Convert status string to RequestStatus
	var status models.RequestStatus
	switch data.Status {
	case "SUBMITTED":
		status = models.StatusSubmitted
	case "APPROVED":
		status = models.StatusApproved
	case "REJECTED":
		status = models.StatusRejected
	case "PAID":
		status = models.StatusPaid
	default:
		status = models.StatusSubmitted // Default fallback
	}

	// Create the request model
	request := &models.Request{
		Title:            data.Title,
		Description:      data.Description,
		Status:           status,
		RejectionReason:  data.RejectionReason,
		ApprovedAt:       data.ApprovedAt,
		RejectedAt:       data.RejectedAt,
		QRIdentifier:     &data.QRIdentifier,
	}

	// Set the base model fields
	request.ID = data.ID
	request.CreatedAt = data.CreatedAt
	request.UpdatedAt = data.UpdatedAt

	// Set serial number
	request.SerialNumber = &data.SerialNumber

	// Set user IDs
	request.CreatedByUserID = &data.CreateByUserID
	if data.ApprovedByUserID != nil {
		request.ApprovedByUserID = data.ApprovedByUserID
	}

	// Set request type ID
	request.RequestTypeID = &data.RequestTypeID

	return request
}
