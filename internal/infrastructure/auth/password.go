package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type Argon2Params struct {
	Memory     uint32
	Time       uint32
	Threads    uint8
	KeyLength  uint32
	SaltLength int
}

var defaultParams = Argon2Params{
	Memory:     64 * 1024,
	Time:       1,
	Threads:    4,
	KeyLength:  32,
	SaltLength: 16,
}

// HashPassword hashes the given password using Argon2id.
func HashPassword(password string) (string, error) {
	salt := make([]byte, defaultParams.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, defaultParams.Time, defaultParams.Memory, defaultParams.Threads, defaultParams.KeyLength)

	// Encode salt and hash to a single string for storage
	encodedHash := fmt.Sprintf("%s.%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hashedPassword),
	)
	return encodedHash, nil
}

// CheckPasswordHash compares a plain password with an Argon2id hashed password.
func CheckPasswordHash(password, encodedHash string) (bool, error) {
	// Split the encoded hash into salt and hash parts
	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %v", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %v", err)
	}

	// Hash the password using the same parameters and salt
	computedHash := argon2.IDKey([]byte(password), salt, defaultParams.Time, defaultParams.Memory, defaultParams.Threads, defaultParams.KeyLength)

	// Compare the hashes
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
