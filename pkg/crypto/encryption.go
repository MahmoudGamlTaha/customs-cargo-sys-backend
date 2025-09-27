package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// EncryptID encrypts an int64 ID using AES encryption and returns a base64 encoded string
func EncryptID(id int64, secretKey string) (string, error) {
	// Convert int64 to bytes
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, uint64(id))

	// Create cipher block
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, idBytes, nil)

	// Encode to base64 URL safe
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptID decrypts a base64 encoded string back to int64 ID
func DecryptID(encryptedID string, secretKey string) (int64, error) {
	// Decode from base64
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedID)
	if err != nil {
		return 0, fmt.Errorf("failed to decode base64: %v", err)
	}

	// Create cipher block
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return 0, fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Check minimum length
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return 0, errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to decrypt: %v", err)
	}

	// Convert bytes back to int64
	if len(plaintext) != 8 {
		return 0, errors.New("invalid plaintext length")
	}

	id := int64(binary.BigEndian.Uint64(plaintext))
	return id, nil
}
