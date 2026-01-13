package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPasswordTooShort indica que la contraseña es muy corta
	ErrPasswordTooShort = errors.New("la contraseña debe tener al menos 6 caracteres")
	// ErrPasswordTooLong indica que la contraseña excede el límite de bcrypt
	ErrPasswordTooLong = errors.New("la contraseña no puede exceder 72 caracteres")
)

const (
	minPasswordLength = 6
	maxPasswordLength = 72 // límite de bcrypt
)

// HashPassword hashea una contraseña usando bcrypt.
// Valida la longitud antes de procesar.
func HashPassword(password string) (string, error) {
	if len(password) < minPasswordLength {
		return "", ErrPasswordTooShort
	}
	if len(password) > maxPasswordLength {
		return "", ErrPasswordTooLong
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compara una contraseña hasheada con una contraseña en texto plano.
// Retorna true si coinciden, false en caso contrario.
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
