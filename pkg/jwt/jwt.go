package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims representa los claims personalizados para el JWT.
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken genera un token JWT con los claims proporcionados.
func GenerateToken(userId uint, email, secret string, expiresIn time.Duration) (string, error) {

	// Crear los claims
	claims := JWTClaims{
		UserID: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	return token.SignedString([]byte(secret))
}

// ValidateToken valida un token JWT y retorna los claims si es válido.
func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	// Parsear token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Obtener claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
