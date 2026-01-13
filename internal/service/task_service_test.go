package service

import (
	"errors"
	"testing"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	userID := uint(1)
	req := &domain.CreateTask{Title: "Nueva tarea"}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Task")).Return(nil)

	task, err := service.Create(userID, req)

	require.NoError(t, err)
	assert.Equal(t, "Nueva tarea", task.Title)
	assert.Equal(t, userID, task.UserID)
	assert.False(t, task.Completed)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_Create_RepoError(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	userID := uint(1)
	req := &domain.CreateTask{Title: "Nueva tarea"}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Task")).Return(errors.New("database error"))

	task, err := service.Create(userID, req)

	assert.Error(t, err)
	assert.Nil(t, task)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	expectedTask := &domain.Task{ID: 1, Title: "Tarea 1", UserID: 1}
	mockRepo.On("GetByID", uint(1)).Return(expectedTask, nil)

	task, err := service.GetByID(1)

	require.NoError(t, err)
	assert.Equal(t, expectedTask.ID, task.ID)
	assert.Equal(t, expectedTask.Title, task.Title)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	mockRepo.On("GetByID", uint(999)).Return(nil, nil)

	task, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Equal(t, ErrTaskNotFound, err)
	assert.Nil(t, task)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_GetByUserID_Success(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	expectedTasks := []domain.Task{
		{ID: 1, Title: "Tarea 1", UserID: 1},
		{ID: 2, Title: "Tarea 2", UserID: 1},
	}
	mockRepo.On("GetByUserID", uint(1)).Return(expectedTasks, nil)

	tasks, err := service.GetByUserID(1)

	require.NoError(t, err)
	assert.Len(t, tasks, 2)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	existingTask := &domain.Task{ID: 1, Title: "Tarea original", UserID: 1}
	newTitle := "Tarea actualizada"
	req := &domain.UpdateTask{Title: &newTitle}

	mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Task")).Return(nil)

	task, err := service.Update(1, 1, req)

	require.NoError(t, err)
	assert.Equal(t, "Tarea actualizada", task.Title)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_Update_NotFound(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	newTitle := "Tarea actualizada"
	req := &domain.UpdateTask{Title: &newTitle}

	mockRepo.On("GetByID", uint(999)).Return(nil, nil)

	task, err := service.Update(999, 1, req)

	assert.Error(t, err)
	assert.Equal(t, ErrTaskNotFound, err)
	assert.Nil(t, task)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_Update_Unauthorized(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	existingTask := &domain.Task{ID: 1, Title: "Tarea", UserID: 1}
	newTitle := "Tarea actualizada"
	req := &domain.UpdateTask{Title: &newTitle}

	mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)

	task, err := service.Update(1, 2, req)

	assert.Error(t, err)
	assert.Equal(t, ErrTaskUnauthorized, err)
	assert.Nil(t, task)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	existingTask := &domain.Task{ID: 1, Title: "Tarea", UserID: 1}

	mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1, 1)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	mockRepo.On("GetByID", uint(999)).Return(nil, nil)

	err := service.Delete(999, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrTaskNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_Delete_Unauthorized(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	existingTask := &domain.Task{ID: 1, Title: "Tarea", UserID: 1}

	mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)

	err := service.Delete(1, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrTaskUnauthorized, err)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_UpdateStatus_Success(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	existingTask := &domain.Task{ID: 1, Title: "Tarea", UserID: 1, Completed: false}

	mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)
	mockRepo.On("UpdateStatus", uint(1), true).Return(nil)

	task, err := service.UpdateStatus(1, 1, true)

	require.NoError(t, err)
	assert.True(t, task.Completed)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_UpdateStatus_Unauthorized(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryMock)
	service := NewTaskService(mockRepo)

	existingTask := &domain.Task{ID: 1, Title: "Tarea", UserID: 1}

	mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)

	task, err := service.UpdateStatus(1, 2, true)

	assert.Error(t, err)
	assert.Equal(t, ErrTaskUnauthorized, err)
	assert.Nil(t, task)
	mockRepo.AssertExpectations(t)
}
