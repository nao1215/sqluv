package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// EncryptPassword encrypts a password string using AES-GCM
func EncryptPassword(password string) (string, error) {
	// Generate a key based on a unique machine identifier
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)

	// Return base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword decrypts an encrypted password string
func DecryptPassword(encryptedPassword string) (string, error) {
	// Get the encryption key
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	// Decode the base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Verify the ciphertext is valid
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// getEncryptionKey generates a stable key based on machine-specific identifiers
func getEncryptionKey() ([]byte, error) {
	// Use the user's home directory path as a unique identifier
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Create a hash of the user's home directory
	hash := sha256.Sum256([]byte(homeDir))
	return hash[:], nil
}

// IsEncrypted checks if a password string is likely encrypted
func IsEncrypted(password string) bool {
	// Check if it looks like a base64 encoded string
	_, err := base64.StdEncoding.DecodeString(password)
	return err == nil && len(password) > 20 // Arbitrary length check
}
