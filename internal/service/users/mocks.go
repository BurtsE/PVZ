package users

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	domain "pvz/internal/domain/users"
)

// MockToken - мок для токен-сервиса
type MockToken struct {
	mock.Mock
}

func (m *MockToken) Create(user domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockToken) Parse(token string) (domain.User, error) {
	args := m.Called(token)
	return args.Get(0).(domain.User), args.Error(1)
}

// MockUserRepo - мок для репозитория
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserById(ctx context.Context, id uuid.UUID) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUser(ctx context.Context, email string) (domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

// Вспомогательные функции для тестов
func validUser() domain.User {
	user, _ := domain.CreateUser("test@example.com", "user", []byte("hash"))
	return user
}
