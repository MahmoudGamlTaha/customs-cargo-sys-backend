package models

import (
	"net"
	"time"

	
)

// UserSession represents a user session entity
type UserSession struct {
	BaseModel
	UserID    int64 `json:"user_id" db:"user_id" validate:"required"`
	TokenHash string    `json:"-" db:"token_hash" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at" validate:"required"`
	LastUsed  time.Time `json:"last_used" db:"last_used"`
	IPAddress net.IP    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`

	// Relationships
	User *User `json:"user,omitempty"`
}

// SerialSequence represents the serial number sequence entity
type SerialSequence struct {
	ID             int       `json:"id" db:"id"`
	Year           int       `json:"year" db:"year"`
	SequenceNumber int       `json:"sequence_number" db:"sequence_number"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// IsExpired checks if the session is expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid (not expired and user is active)
func (s *UserSession) IsValid() bool {
	if s.IsExpired() {
		return false
	}
	if s.User != nil && !s.User.IsActive {
		return false
	}
	return true
}

// UpdateLastUsed updates the last used timestamp
func (s *UserSession) UpdateLastUsed() {
	s.LastUsed = time.Now()
}

// SessionResponse represents the response payload for session data
type SessionResponse struct {
	ID        int64 `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
	LastUsed  time.Time `json:"last_used"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts UserSession to SessionResponse
func (s *UserSession) ToResponse() *SessionResponse {
	return &SessionResponse{
		ID:        s.ID,
		ExpiresAt: s.ExpiresAt,
		LastUsed:  s.LastUsed,
		IPAddress: s.IPAddress.String(),
		UserAgent: s.UserAgent,
		CreatedAt: s.CreatedAt,
	}
}

