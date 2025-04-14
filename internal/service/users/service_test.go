package users

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	domain "pvz/internal/domain/users"
	"testing"
)

func TestUserService_RegisterUser(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		ctx := context.Background()
		mockToken := NewToken([]byte("mock-token"))
		mockRepo := new(MockUserRepo)
		service := &UserService{
			token:    mockToken,
			userRepo: mockRepo,
		}

		// Ожидаемые вызовы
		//testUser := validUser()
		mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil)

		// Вызов метода
		token, err := service.RegisterUser(ctx, "test@example.com", "validPass123", domain.Employee)

		// Проверки
		require.NoError(t, err)
		user, err := mockToken.Parse(token)
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", user.Email())
		assert.Equal(t, domain.Employee, user.Role())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		service := &UserService{}
		_, err := service.RegisterUser(context.Background(), "test@example.com", "short", "employee")
		assert.Error(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		mockToken := NewToken([]byte("mock-token"))
		mockRepo := new(MockUserRepo)
		service := &UserService{
			token:    mockToken,
			userRepo: mockRepo,
		}

		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("users.User")).Return(errors.New("db error"))

		_, err := service.RegisterUser(ctx, "test@example.com", "validPass123", "employee")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_ValidateUser(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		mockToken := NewToken([]byte("mock-token"))
		service := &UserService{token: mockToken}

		expectedUser, err := domain.CreateUser("test@example.com", domain.Employee, []byte("hash"))
		require.NoError(t, err)
		token, err := mockToken.Create(expectedUser)
		require.NoError(t, err)
		user, err := service.ValidateUser(context.Background(), token)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("invalid token", func(t *testing.T) {
		mockToken := NewToken([]byte("mock-token"))
		service := &UserService{token: mockToken}

		_, err := service.ValidateUser(context.Background(), "invalid-token")
		assert.Error(t, err)
	})
}

func TestUserService_LoginUser(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		ctx := context.Background()
		mockToken := NewToken([]byte("mock-token-sfdd"))
		mockRepo := new(MockUserRepo)
		service := &UserService{
			token:    mockToken,
			userRepo: mockRepo,
		}
		hash, err := hashPassword("password1")
		require.NoError(t, err)

		testUser, err := domain.CreateUser("test@example.com", domain.Moderator, hash)
		require.NoError(t, err)
		mockRepo.On("GetUser", ctx, "test@example.com").Return(testUser, nil)

		tokenExpected, err := mockToken.Create(testUser)
		require.NoError(t, err)
		// Подменяем функцию проверки пароля
		//oldCompare := compareHashAndPassword
		//compareHashAndPassword = func(hash []byte, password string) error { return nil }
		//defer func() { compareHashAndPassword = oldCompare }()

		token, err := service.LoginUser(ctx, "test@example.com", "password1")
		require.NoError(t, err)
		assert.Equal(t, tokenExpected, token)
	})

	t.Run("wrong password", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUserRepo)
		service := &UserService{userRepo: mockRepo}

		testUser := validUser()
		mockRepo.On("GetUser", ctx, "test@example.com").Return(testUser, nil)

		// Подменяем функцию проверки пароля
		//oldCompare := compareHashAndPassword
		//compareHashAndPassword = func(hash []byte, password string) error { return errors.New("mismatch") }
		//defer func() { compareHashAndPassword = oldCompare }()

		_, err := service.LoginUser(ctx, "test@example.com", "wrong-pass")
		assert.Error(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUserRepo)
		service := &UserService{userRepo: mockRepo}

		mockRepo.On("GetUser", ctx, "nonexistent@test.com").Return(domain.User{}, errors.New("not found"))

		_, err := service.LoginUser(ctx, "nonexistent@test.com", "pass")
		assert.Error(t, err)
	})
}
