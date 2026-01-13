package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexroel/gin-tasks-api/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token no proporcionado")
}

func TestAuthMiddleware_InvalidFormat_NoBearer(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidToken123")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Formato de token inválido")
}

func TestAuthMiddleware_InvalidFormat_WrongPrefix(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic sometoken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Formato de token inválido")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token inválido")
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Generar un token que ya expiró
	expiredToken, err := jwt.GenerateToken(1, "test@example.com", "test-secret", -1*time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token inválido")
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	// Generar token con un secret diferente
	token, err := jwt.GenerateToken(1, "test@example.com", "different-secret", 24*time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token inválido")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "test-secret"
	userID := uint(123)
	email := "test@example.com"

	token, err := jwt.GenerateToken(userID, email, secret, 24*time.Hour)
	assert.NoError(t, err)

	var capturedUserID uint
	var capturedEmail string

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) {
		capturedUserID = c.GetUint("userID")
		capturedEmail = c.GetString("userEmail")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userID, capturedUserID)
	assert.Equal(t, email, capturedEmail)
}

func TestAuthMiddleware_ValidToken_SetsContext(t *testing.T) {
	secret := "my-secret-key"
	userID := uint(456)
	email := "user@test.com"

	token, err := jwt.GenerateToken(userID, email, secret, 1*time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) {
		// Verificar que los valores están en el contexto
		id, exists := c.Get("userID")
		assert.True(t, exists)
		assert.Equal(t, userID, id)

		mail, exists := c.Get("userEmail")
		assert.True(t, exists)
		assert.Equal(t, email, mail)

		c.JSON(http.StatusOK, gin.H{"userID": id, "email": mail})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Tests para GetUserID helper
func TestGetUserID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedID := uint(99)
	c.Set("userID", expectedID)

	userID, exists := GetUserID(c)

	assert.True(t, exists)
	assert.Equal(t, expectedID, userID)
}

func TestGetUserID_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID, exists := GetUserID(c)

	assert.False(t, exists)
	assert.Equal(t, uint(0), userID)
}

func TestAuthMiddleware_EmptyBearerToken(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_MultipleSpaces(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer  token  extra")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
