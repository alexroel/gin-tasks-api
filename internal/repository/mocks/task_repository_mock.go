package mocks

import (
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

// TaskRepositoryMock es un mock del repositorio de tareas
type TaskRepositoryMock struct {
	mock.Mock
}

func (m *TaskRepositoryMock) Create(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *TaskRepositoryMock) GetAll() ([]domain.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *TaskRepositoryMock) GetByID(id uint) (*domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *TaskRepositoryMock) GetByUserID(userID uint) ([]domain.Task, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *TaskRepositoryMock) Update(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *TaskRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *TaskRepositoryMock) UpdateStatus(id uint, completed bool) error {
	args := m.Called(id, completed)
	return args.Error(0)
}
