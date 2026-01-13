package mocks

import (
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

// TaskServiceMock es un mock del servicio de tareas
type TaskServiceMock struct {
	mock.Mock
}

func (m *TaskServiceMock) Create(userID uint, req *domain.CreateTask) (*domain.Task, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *TaskServiceMock) GetAll() ([]domain.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *TaskServiceMock) GetByID(id uint) (*domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *TaskServiceMock) GetByUserID(userID uint) ([]domain.Task, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *TaskServiceMock) Update(id, userID uint, req *domain.UpdateTask) (*domain.Task, error) {
	args := m.Called(id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *TaskServiceMock) Delete(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *TaskServiceMock) UpdateStatus(id, userID uint, completed bool) (*domain.Task, error) {
	args := m.Called(id, userID, completed)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}
