package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordConfig represents password hashing configuration
type PasswordConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultPasswordConfig returns default password configuration
func DefaultPasswordConfig() *PasswordConfig {
	return &PasswordConfig{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// PasswordManager handles password operations
type PasswordManager struct {
	config *PasswordConfig
}

// NewPasswordManager creates a new password manager
func NewPasswordManager(config *PasswordConfig) *PasswordManager {
	if config == nil {
		config = DefaultPasswordConfig()
	}
	return &PasswordManager{config: config}
}

// HashPassword hashes a password using Argon2id
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	// Generate a random salt
	salt, err := generateRandomBytes(pm.config.SaltLength)
	if err != nil {
		return "", err
	}

	// Hash the password
	hash := argon2.IDKey([]byte(password), salt, pm.config.Iterations, pm.config.Memory, pm.config.Parallelism, pm.config.KeyLength)

	// Encode the salt and hash to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return the encoded hash in the format: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, pm.config.Memory, pm.config.Iterations, pm.config.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword verifies a password against a hash
func (pm *PasswordManager) VerifyPassword(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, errors.New("incompatible version of argon2")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters
	otherHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	// Check that the contents of the hashed passwords are identical
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// generateRandomBytes generates random bytes of the specified length
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if len(password) > 100 {
		return errors.New("password must be less than 100 characters long")
	}
	return nil
}

