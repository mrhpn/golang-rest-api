package security

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/mrhpn/go-rest-api/internal/constants"
)

// HashPassword hashes a raw string using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.DefaultBcryptCost)
	return string(bytes), err
}

// CheckPassword compares a hash with a raw password
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
