package utils

import (
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength is the minimum password length
	MinPasswordLength = 8
	// MaxPasswordLength is the maximum password length
	MaxPasswordLength = 72
	// BcryptCost is the cost factor for bcrypt hashing
	BcryptCost = 12
)

// PasswordStrength represents password strength level
type PasswordStrength int

const (
	PasswordWeak PasswordStrength = iota
	PasswordModerate
	PasswordStrong
	PasswordVeryStrong
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	if err := ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// ComparePassword compares a hashed password with a plain text password
func ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password must not exceed %d characters", MaxPasswordLength)
	}

	strength := CheckPasswordStrength(password)

	// Only accept Strong & VeryStrong
	if strength == PasswordWeak || strength == PasswordModerate {
		return fmt.Errorf("password too weak or moderate, must be strong or very strong")
	}

	return nil
}

// CheckPasswordStrength checks the strength of a password
func CheckPasswordStrength(password string) PasswordStrength {
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	specialCount := 0

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
			specialCount++
		}
	}

	length := len(password)

	// Very Strong: ≥12 + full requirements + ≥2 special chars
	if length >= 12 && hasUpper && hasLower && hasNumber && hasSpecial && specialCount >= 2 {
		return PasswordVeryStrong
	}

	// Strong: ≥10 + all requirements
	if length >= 10 && hasUpper && hasLower && hasNumber && hasSpecial {
		return PasswordStrong
	}

	// Moderate: ≥8 + uppercase + lowercase + number, but NO special char
	if length >= 8 && hasUpper && hasLower && hasNumber && !hasSpecial {
		return PasswordModerate
	}

	// Weak
	return PasswordWeak
}

// GetPasswordStrengthText returns human-readable password strength
func GetPasswordStrengthText(password string) string {
	switch CheckPasswordStrength(password) {
	case PasswordVeryStrong:
		return "Very Strong"
	case PasswordStrong:
		return "Strong"
	case PasswordModerate:
		return "Moderate"
	default:
		return "Weak"
	}
}
