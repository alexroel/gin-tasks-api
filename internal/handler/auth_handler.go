package handler

import (
	"net/http"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/middleware"
	"github.com/alexroel/gin-tasks-api/internal/service"
	"github.com/alexroel/gin-tasks-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// SignUpHandler godoc
// @Summary      Registro de usuario
// @Description  Registra un nuevo usuario en el sistema
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      domain.UserCreate  true  "Datos de registro del usuario"
// @Success      201   {object}  utils.Response{data=domain.UserResponse} "Usuario registrado"
// @Failure      400   {object}  utils.Response "Datos inválidos"
// @Router       /auth/signup [post]
func (h *AuthHandler) SignUpHandler(c *gin.Context) {
	var req domain.UserCreate

	// Validar vody
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos invávlidos"+err.Error())
		return
	}

	// Registrar usuario
	user, err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Responder con el usuario creado
	utils.SuccessResponse(c, http.StatusCreated, "Usuario registrado exitosamente", gin.H{
		"user": gin.H{
			"id":        user.ID,
			"full_name": user.FullName,
			"email":     user.Email,
		},
	})

}

// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica un usuario y retorna un token JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body domain.UserLogin true "Credenciales del usuario"
// @Success      200 {object} utils.Response{data=object{token=string,user=domain.UserResponse}} "Login exitoso"
// @Failure      400 {object} utils.Response "Datos inválidos"
// @Failure      401 {object} utils.Response "Credenciales incorrectas"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.UserLogin

	// Validar request body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Autenticar usuario
	token, user, err := h.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Inicio de sesión exitoso", gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"full_name": user.FullName,
			"email":     user.Email,
		},
	})
}

// Profile godoc
// @Summary      Obtener perfil
// @Description  Obtiene la información del usuario autenticado
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=domain.UserResponse} "Perfil obtenido"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      404 {object} utils.Response "Usuario no encontrado"
// @Router       /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	// Obtener ID del usuario del contexto
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Obtener usuario
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil obtenido exitosamente", user.ToResponse())
}

// UpdateProfile godoc
// @Summary      Actualizar perfil
// @Description  Actualiza la información del usuario autenticado
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.UserUpdate true "Datos a actualizar"
// @Success      200 {object} utils.Response{data=domain.UserResponse} "Perfil actualizado"
// @Failure      400 {object} utils.Response "Datos inválidos"
// @Failure      401 {object} utils.Response "No autenticado"
// @Router       /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Obtener ID del usuario del contexto
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	var req domain.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil actualizado exitosamente", user.ToResponse())
}

// DeleteAccount godoc
// @Summary      Eliminar cuenta
// @Description  Elimina la cuenta del usuario autenticado
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response "Cuenta eliminada"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      500 {object} utils.Response "Error interno"
// @Router       /auth/profile [delete]
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	// Obtener ID del usuario del contexto
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	if err := h.authService.DeleteAccount(userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar la cuenta: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cuenta eliminada exitosamente", nil)
}
