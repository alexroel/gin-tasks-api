package handler

import (
	"net/http"
	"strconv"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/middleware"
	"github.com/alexroel/gin-tasks-api/internal/service"
	"github.com/alexroel/gin-tasks-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskService
}

// NewTaskHandler crea una nueva instancia de TaskHandler
func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// Create godoc
// @Summary      Crear tarea
// @Description  Crea una nueva tarea para el usuario autenticado
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateTask true "Datos de la tarea"
// @Success      201 {object} utils.Response{data=domain.TaskResponse} "Tarea creada"
// @Failure      400 {object} utils.Response "Datos inválidos"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      500 {object} utils.Response "Error interno"
// @Router       /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	// Obtener ID del usuario autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	var req domain.CreateTask
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	task, err := h.taskService.Create(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al crear la tarea: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Tarea creada exitosamente", task.ToResponse())
}

// GetAll godoc
// @Summary      Listar tareas
// @Description  Obtiene todas las tareas del usuario autenticado
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=[]domain.TaskResponse} "Lista de tareas"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      500 {object} utils.Response "Error interno"
// @Router       /tasks [get]
func (h *TaskHandler) GetAll(c *gin.Context) {
	// Obtener ID del usuario autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	tasks, err := h.taskService.GetByUserID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener las tareas: "+err.Error())
		return
	}

	// Convertir a respuesta
	var tasksResponse []domain.TaskResponse
	for _, task := range tasks {
		tasksResponse = append(tasksResponse, task.ToResponse())
	}

	utils.SuccessResponse(c, http.StatusOK, "Tareas obtenidas exitosamente", tasksResponse)
}

// GetByID godoc
// @Summary      Obtener tarea
// @Description  Obtiene una tarea por su ID
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la tarea"
// @Success      200 {object} utils.Response{data=domain.TaskResponse} "Tarea obtenida"
// @Failure      400 {object} utils.Response "ID inválido"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      403 {object} utils.Response "Sin permiso"
// @Failure      404 {object} utils.Response "Tarea no encontrada"
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	// Obtener ID del usuario autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Obtener ID de la tarea
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de tarea inválido")
		return
	}

	task, err := h.taskService.GetByID(uint(taskID))
	if err != nil {
		if err == service.ErrTaskNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Tarea no encontrada")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener la tarea: "+err.Error())
		return
	}

	// Verificar que la tarea pertenece al usuario
	if task.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "No tienes permiso para ver esta tarea")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tarea obtenida exitosamente", task.ToResponse())
}

// Update godoc
// @Summary      Actualizar tarea
// @Description  Actualiza una tarea existente
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la tarea"
// @Param        request body domain.UpdateTask true "Datos a actualizar"
// @Success      200 {object} utils.Response{data=domain.TaskResponse} "Tarea actualizada"
// @Failure      400 {object} utils.Response "Datos inválidos"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      403 {object} utils.Response "Sin permiso"
// @Failure      404 {object} utils.Response "Tarea no encontrada"
// @Router       /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	// Obtener ID del usuario autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Obtener ID de la tarea
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de tarea inválido")
		return
	}

	var req domain.UpdateTask
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	task, err := h.taskService.Update(uint(taskID), userID, &req)
	if err != nil {
		switch err {
		case service.ErrTaskNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Tarea no encontrada")
		case service.ErrTaskUnauthorized:
			utils.ErrorResponse(c, http.StatusForbidden, "No tienes permiso para modificar esta tarea")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar la tarea: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tarea actualizada exitosamente", task.ToResponse())
}

// Delete godoc
// @Summary      Eliminar tarea
// @Description  Elimina una tarea por su ID
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la tarea"
// @Success      200 {object} utils.Response "Tarea eliminada"
// @Failure      400 {object} utils.Response "ID inválido"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      403 {object} utils.Response "Sin permiso"
// @Failure      404 {object} utils.Response "Tarea no encontrada"
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	// Obtener ID del usuario autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Obtener ID de la tarea
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de tarea inválido")
		return
	}

	err = h.taskService.Delete(uint(taskID), userID)
	if err != nil {
		switch err {
		case service.ErrTaskNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Tarea no encontrada")
		case service.ErrTaskUnauthorized:
			utils.ErrorResponse(c, http.StatusForbidden, "No tienes permiso para eliminar esta tarea")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar la tarea: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tarea eliminada exitosamente", nil)
}

// ToggleStatus godoc
// @Summary      Cambiar estado de tarea
// @Description  Actualiza el estado de completado de una tarea
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la tarea"
// @Param        request body object{completed=bool} true "Nuevo estado"
// @Success      200 {object} utils.Response{data=domain.TaskResponse} "Estado actualizado"
// @Failure      400 {object} utils.Response "Datos inválidos"
// @Failure      401 {object} utils.Response "No autenticado"
// @Failure      403 {object} utils.Response "Sin permiso"
// @Failure      404 {object} utils.Response "Tarea no encontrada"
// @Router       /tasks/{id}/status [patch]
func (h *TaskHandler) ToggleStatus(c *gin.Context) {
	// Obtener ID del usuario autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Obtener ID de la tarea
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de tarea inválido")
		return
	}

	// Estructura para recibir el estado
	var req struct {
		Completed bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	task, err := h.taskService.UpdateStatus(uint(taskID), userID, req.Completed)
	if err != nil {
		switch err {
		case service.ErrTaskNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Tarea no encontrada")
		case service.ErrTaskUnauthorized:
			utils.ErrorResponse(c, http.StatusForbidden, "No tienes permiso para modificar esta tarea")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar el estado: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Estado de tarea actualizado exitosamente", task.ToResponse())
}
