package mocks

import (
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

// AuthServiceMock es un mock del servicio de autenticación
type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Register(req *domain.UserCreate) (*domain.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *AuthServiceMock) Login(req *domain.UserLogin) (string, *domain.User, error) {
	args := m.Called(req)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*domain.User), args.Error(2)
}

func (m *AuthServiceMock) GetUserByID(userID uint) (*domain.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *AuthServiceMock) UpdateProfile(userID uint, req *domain.UserUpdate) (*domain.User, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *AuthServiceMock) DeleteAccount(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
