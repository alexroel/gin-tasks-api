package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-12345"

func TestGenerateToken_Success(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	expiresIn := time.Hour

	token, err := GenerateToken(userID, email, testSecret, expiresIn)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_Success(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	expiresIn := time.Hour

	token, _ := GenerateToken(userID, email, testSecret, expiresIn)

	claims, err := ValidateToken(token, testSecret)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	expiresIn := time.Hour

	token, _ := GenerateToken(userID, email, testSecret, expiresIn)

	claims, err := ValidateToken(token, "wrong-secret")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	// Token que expira inmediatamente
	expiresIn := -time.Hour

	token, _ := GenerateToken(userID, email, testSecret, expiresIn)

	claims, err := ValidateToken(token, testSecret)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	claims, err := ValidateToken("invalid-token", testSecret)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	claims, err := ValidateToken("", testSecret)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerateToken_DifferentUsers(t *testing.T) {
	expiresIn := time.Hour

	token1, _ := GenerateToken(1, "user1@example.com", testSecret, expiresIn)
	token2, _ := GenerateToken(2, "user2@example.com", testSecret, expiresIn)

	// Tokens diferentes
	assert.NotEqual(t, token1, token2)

	// Validar claims correctos
	claims1, _ := ValidateToken(token1, testSecret)
	claims2, _ := ValidateToken(token2, testSecret)

	assert.Equal(t, uint(1), claims1.UserID)
	assert.Equal(t, uint(2), claims2.UserID)
}

// ========== Table-Driven Tests ==========

func TestGenerateToken_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		email     string
		secret    string
		expiresIn time.Duration
		wantErr   bool
	}{
		{
			name:      "token válido - usuario normal",
			userID:    1,
			email:     "user@example.com",
			secret:    testSecret,
			expiresIn: time.Hour,
			wantErr:   false,
		},
		{
			name:      "token válido - ID grande",
			userID:    999999,
			email:     "bigid@example.com",
			secret:    testSecret,
			expiresIn: time.Hour,
			wantErr:   false,
		},
		{
			name:      "token válido - duración larga",
			userID:    1,
			email:     "long@example.com",
			secret:    testSecret,
			expiresIn: 24 * 7 * time.Hour, // 1 semana
			wantErr:   false,
		},
		{
			name:      "token válido - duración corta",
			userID:    1,
			email:     "short@example.com",
			secret:    testSecret,
			expiresIn: time.Minute,
			wantErr:   false,
		},
		{
			name:      "token válido - secret largo",
			userID:    1,
			email:     "longsecret@example.com",
			secret:    "this-is-a-very-long-secret-key-that-should-work-fine-12345",
			expiresIn: time.Hour,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email, tt.secret, tt.expiresIn)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verificar que el token se puede validar
				claims, err := ValidateToken(token, tt.secret)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.email, claims.Email)
			}
		})
	}
}

func TestValidateToken_TableDriven(t *testing.T) {
	// Generar tokens para las pruebas
	validToken, _ := GenerateToken(1, "test@example.com", testSecret, time.Hour)
	expiredToken, _ := GenerateToken(1, "test@example.com", testSecret, -time.Hour)
	differentSecretToken, _ := GenerateToken(1, "test@example.com", "different-secret", time.Hour)

	tests := []struct {
		name           string
		token          string
		secret         string
		wantErr        bool
		expectedUserID uint
		expectedEmail  string
	}{
		{
			name:           "token válido",
			token:          validToken,
			secret:         testSecret,
			wantErr:        false,
			expectedUserID: 1,
			expectedEmail:  "test@example.com",
		},
		{
			name:    "token expirado",
			token:   expiredToken,
			secret:  testSecret,
			wantErr: true,
		},
		{
			name:    "secret incorrecto",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
		},
		{
			name:    "token generado con otro secret",
			token:   differentSecretToken,
			secret:  testSecret,
			wantErr: true,
		},
		{
			name:    "token vacío",
			token:   "",
			secret:  testSecret,
			wantErr: true,
		},
		{
			name:    "token inválido - texto aleatorio",
			token:   "not.a.valid.jwt.token",
			secret:  testSecret,
			wantErr: true,
		},
		{
			name:    "token malformado",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
			secret:  testSecret,
			wantErr: true,
		},
		{
			name:    "secret vacío",
			token:   validToken,
			secret:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, tt.secret)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.expectedUserID, claims.UserID)
				assert.Equal(t, tt.expectedEmail, claims.Email)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	tests := []struct {
		name      string
		expiresIn time.Duration
		shouldExp bool
	}{
		{"expira en 1 hora", time.Hour, false},
		{"expira en 1 minuto", time.Minute, false},
		{"expiró hace 1 hora", -time.Hour, true},
		{"expiró hace 1 segundo", -time.Second, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _ := GenerateToken(1, "test@example.com", testSecret, tt.expiresIn)
			claims, err := ValidateToken(token, testSecret)

			if tt.shouldExp {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}
