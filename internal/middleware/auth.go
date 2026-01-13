package middleware

import (
	"net/http"
	"strings"

	"github.com/alexroel/gin-tasks-api/pkg/jwt"
	"github.com/alexroel/gin-tasks-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del encabezado Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token no proporcionado")
			c.Abort() // Detener la ejecución del middleware
			return
		}

		// Verificar formato del token "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Formato de token inválido")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token inválido: "+err.Error())
			c.Abort()
			return
		}

		// Almacenar los claims en el contexto para su uso posterior
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

		c.Next() // Continuar con la siguiente función en la cadena de middleware
	}
}

// GetUserID obtiene el ID del usuario del contexto
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}
