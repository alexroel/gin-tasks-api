package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_Success(t *testing.T) {
	password := "password123"

	hash, err := HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestHashPassword_TooShort(t *testing.T) {
	password := "12345" // 5 caracteres, mínimo es 6

	hash, err := HashPassword(password)

	assert.Error(t, err)
	assert.Equal(t, ErrPasswordTooShort, err)
	assert.Empty(t, hash)
}

func TestHashPassword_TooLong(t *testing.T) {
	// Crear contraseña de 73 caracteres (máximo es 72)
	password := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	hash, err := HashPassword(password)

	assert.Error(t, err)
	assert.Equal(t, ErrPasswordTooLong, err)
	assert.Empty(t, hash)
}

func TestCheckPassword_Valid(t *testing.T) {
	password := "password123"
	hash, _ := HashPassword(password)

	result := CheckPassword(hash, password)

	assert.True(t, result)
}

func TestCheckPassword_Invalid(t *testing.T) {
	password := "password123"
	hash, _ := HashPassword(password)

	result := CheckPassword(hash, "wrongpassword")

	assert.False(t, result)
}

func TestCheckPassword_EmptyHash(t *testing.T) {
	result := CheckPassword("", "password123")

	assert.False(t, result)
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "password123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// bcrypt genera hashes diferentes cada vez
	assert.NotEqual(t, hash1, hash2)

	// Pero ambos deben validar correctamente
	assert.True(t, CheckPassword(hash1, password))
	assert.True(t, CheckPassword(hash2, password))
}

// ========== Table-Driven Tests ==========

func TestHashPassword_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "contraseña válida mínima",
			password:    "123456",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "contraseña válida normal",
			password:    "miContraseñaSegura123",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "contraseña válida con caracteres especiales",
			password:    "P@ssw0rd!#$%",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "contraseña válida máxima (72 chars)",
			password:    strings.Repeat("a", 72),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "contraseña muy corta - 5 chars",
			password:    "12345",
			wantErr:     true,
			expectedErr: ErrPasswordTooShort,
		},
		{
			name:        "contraseña muy corta - 1 char",
			password:    "a",
			wantErr:     true,
			expectedErr: ErrPasswordTooShort,
		},
		{
			name:        "contraseña muy corta - vacía",
			password:    "",
			wantErr:     true,
			expectedErr: ErrPasswordTooShort,
		},
		{
			name:        "contraseña muy larga - 73 chars",
			password:    strings.Repeat("a", 73),
			wantErr:     true,
			expectedErr: ErrPasswordTooLong,
		},
		{
			name:        "contraseña muy larga - 100 chars",
			password:    strings.Repeat("x", 100),
			wantErr:     true,
			expectedErr: ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				// Verificar que el hash funciona
				assert.True(t, CheckPassword(hash, tt.password))
			}
		})
	}
}

func TestCheckPassword_TableDriven(t *testing.T) {
	// Generar un hash válido para las pruebas
	validPassword := "validPassword123"
	validHash, _ := HashPassword(validPassword)

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		{
			name:     "contraseña correcta",
			hash:     validHash,
			password: validPassword,
			want:     true,
		},
		{
			name:     "contraseña incorrecta",
			hash:     validHash,
			password: "wrongPassword",
			want:     false,
		},
		{
			name:     "contraseña con espacio extra",
			hash:     validHash,
			password: validPassword + " ",
			want:     false,
		},
		{
			name:     "contraseña vacía",
			hash:     validHash,
			password: "",
			want:     false,
		},
		{
			name:     "hash vacío",
			hash:     "",
			password: validPassword,
			want:     false,
		},
		{
			name:     "hash inválido",
			hash:     "invalid-hash",
			password: validPassword,
			want:     false,
		},
		{
			name:     "hash y contraseña vacíos",
			hash:     "",
			password: "",
			want:     false,
		},
		{
			name:     "contraseña case sensitive - mayúsculas",
			hash:     validHash,
			password: strings.ToUpper(validPassword),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.hash, tt.password)
			assert.Equal(t, tt.want, result)
		})
	}
}
