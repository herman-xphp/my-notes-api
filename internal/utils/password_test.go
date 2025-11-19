package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "SecurePass123!"

	hashed, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestHashPasswordTooShort(t *testing.T) {
	password := "Short1!"

	_, err := HashPassword(password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least")
}

func TestHashPasswordTooWeak(t *testing.T) {
	password := "weakpassword"

	_, err := HashPassword(password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too weak")
}

func TestComparePasswordSuccess(t *testing.T) {
	password := "SecurePass123!"

	hashed, err := HashPassword(password)
	assert.NoError(t, err)

	err = ComparePassword(hashed, password)
	assert.NoError(t, err)
}

func TestComparePasswordFailed(t *testing.T) {
	password := "SecurePass123!"
	wrongPassword := "WrongPass123!"

	hashed, err := HashPassword(password)
	assert.NoError(t, err)

	err = ComparePassword(hashed, wrongPassword)
	assert.Error(t, err)
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid strong password",
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:     "too short",
			password: "Short1!",
			wantErr:  true,
		},
		{
			name:     "no uppercase",
			password: "securepass123!",
			wantErr:  true,
		},
		{
			name:     "no lowercase",
			password: "SECUREPASS123!",
			wantErr:  true,
		},
		{
			name:     "no number",
			password: "SecurePass!",
			wantErr:  true,
		},
		{
			name:     "no special char",
			password: "SecurePass123",
			wantErr:  true,
		},
		{
			name:     "moderate password",
			password: "Password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     PasswordStrength
	}{
		{
			name:     "very strong password",
			password: "VerySecure123!@#",
			want:     PasswordVeryStrong,
		},
		{
			name:     "strong password",
			password: "SecurePass123!",
			want:     PasswordStrong,
		},
		{
			name:     "moderate password",
			password: "Password123",
			want:     PasswordModerate,
		},
		{
			name:     "weak password",
			password: "password",
			want:     PasswordWeak,
		},
		{
			name:     "weak password short",
			password: "Pass1!",
			want:     PasswordWeak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordStrength(tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPasswordStrengthText(t *testing.T) {
	tests := []struct {
		password string
		want     string
	}{
		{"VerySecure123!@#", "Very Strong"},
		{"SecurePass123!", "Strong"},
		{"Password123", "Moderate"},
		{"password", "Weak"},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			got := GetPasswordStrengthText(tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}
