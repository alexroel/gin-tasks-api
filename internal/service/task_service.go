package service

import (
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/repository"
)

var (
	ErrTaskNotFound    = errors.New("tarea no encontrada")
	ErrTaskUnauthorized = errors.New("no tienes permiso para acceder a esta tarea")
)

// TaskService define las operaciones de negocio para tareas
type TaskService interface {
	Create(userID uint, req *domain.CreateTask) (*domain.Task, error)
	GetAll() ([]domain.Task, error)
	GetByID(id uint) (*domain.Task, error)
	GetByUserID(userID uint) ([]domain.Task, error)
	Update(id, userID uint, req *domain.UpdateTask) (*domain.Task, error)
	Delete(id, userID uint) error
	UpdateStatus(id, userID uint, completed bool) (*domain.Task, error)
}

type taskService struct {
	repo repository.TaskRepository
}

// NewTaskService crea una nueva instancia de TaskService
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// Create crea una nueva tarea para un usuario
func (s *taskService) Create(userID uint, req *domain.CreateTask) (*domain.Task, error) {
	task := &domain.Task{
		Title:  req.Title,
		UserID: userID,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetAll obtiene todas las tareas
func (s *taskService) GetAll() ([]domain.Task, error) {
	return s.repo.GetAll()
}

// GetByID obtiene una tarea por su ID
func (s *taskService) GetByID(id uint) (*domain.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// GetByUserID obtiene todas las tareas de un usuario
func (s *taskService) GetByUserID(userID uint) ([]domain.Task, error) {
	return s.repo.GetByUserID(userID)
}

// Update actualiza una tarea existente
func (s *taskService) Update(id, userID uint, req *domain.UpdateTask) (*domain.Task, error) {
	// Obtener tarea existente
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Verificar que la tarea pertenece al usuario
	if task.UserID != userID {
		return nil, ErrTaskUnauthorized
	}

	// Actualizar campos si se proporcionan
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

// Delete elimina una tarea por su ID
func (s *taskService) Delete(id, userID uint) error {
	// Obtener tarea existente
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	// Verificar que la tarea pertenece al usuario
	if task.UserID != userID {
		return ErrTaskUnauthorized
	}

	return s.repo.Delete(id)
}

// UpdateStatus actualiza el estado de completado de una tarea
func (s *taskService) UpdateStatus(id, userID uint, completed bool) (*domain.Task, error) {
	// Obtener tarea existente
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Verificar que la tarea pertenece al usuario
	if task.UserID != userID {
		return nil, ErrTaskUnauthorized
	}

	// Actualizar estado
	if err := s.repo.UpdateStatus(id, completed); err != nil {
		return nil, err
	}

	task.Completed = completed
	return task, nil
}
