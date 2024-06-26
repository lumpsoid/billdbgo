package parser

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
)

func getPasswordHash(password string) ([]byte, error) {
	passwordBytes := []byte(password)

	hasher := sha256.New()
	hasher.Write(passwordBytes)
	hash := hasher.Sum(nil)

	return hash, nil
}

func decrypt(cyphertext []byte, key []byte) ([]byte, error) {
	blockSize := 12
	// Initialize AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Check if ciphertext is valid
	if len(cyphertext) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	c := len(cyphertext)
	// Extract IV (last 12 bytes of ciphertext)
	iv := cyphertext[c-blockSize : c]

	// Get actual encrypted data (exclude IV)
	data := cyphertext[:c-blockSize]

	// Create AES-GCM cipher instance
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt using AES-GCM
	plaintext, err := aesgcm.Open(nil, iv, data, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
