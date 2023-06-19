package auth

import (
	"crypto/rand"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the input password
func HashPassword(password string, salt []byte) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Get the bcrypt hashed password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 16)
	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), err
}

// CheckPasswordHash checks if hashed password matches raw password
func CheckPasswordHash(hashedPassword, password string, salt []byte) error {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), passwordBytes)
}

// Santinize removes redundant spaces and encodes some special characters which helps to avoid SQL injection
func Santinize(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}

// Generate 16 bytes randomly and securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func GenerateRandomSalt() ([]byte, error) {
	salt := make([]byte, 16)

	_, err := rand.Read(salt[:])
	if err != nil {
		return []byte("error"), err
	}

	return salt, nil
}
