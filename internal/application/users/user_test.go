package users

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz/internal/domain/users"
)

// MockUserService implements UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(ctx context.Context, email, password, role string) (string, error) {
	args := m.Called(ctx, email, password, role)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) ValidateUser(ctx context.Context, token string) (users.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestUserHandlers_PostDummyLogin(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		role         string
		mockSetup    func(*MockUserService)
		expectedCode int
		expectedBody string
	}{
		{
			name:        "successful dummy login",
			requestBody: `{"role": "employee"}`,
			role:        users.Employee,
			mockSetup: func(m *MockUserService) {
				m.On("RegisterUser", mock.Anything, "dummy@example.com", "dummy_password", "employee").
					Return("test-token", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"token":"test-token"}`,
		},
		{
			name:        "invalid role",
			requestBody: `{"role": "invalid"}`,
			role:        "",
			mockSetup: func(m *MockUserService) {
				m.On("RegisterUser", mock.Anything, "dummy@example.com", "dummy_password", "invalid").
					Return("", errors.New("service error"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"service error"}`,
		},
		{
			name:        "service error",
			requestBody: `{"role": "employee"}`,
			role:        users.Employee,
			mockSetup: func(m *MockUserService) {
				m.On("RegisterUser", mock.Anything, "dummy@example.com", "dummy_password", "employee").
					Return("", errors.New("registration failed"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"registration failed"}`,
		},
		{
			name:        "invalid JSON",
			requestBody: `invalid json`,
			mockSetup: func(m *MockUserService) {
				// No mock setup needed as it should fail before service call
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockService := new(MockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}

			// Create handler
			handler := RegisterUserHandlers(mockService)

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/dummy-login", strings.NewReader(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.PostDummyLogin(c)

			// Verify
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandlers_PostLogin(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*MockUserService)
		expectedCode int
		expectedBody string
	}{
		{
			name:        "successful login",
			requestBody: `{"email": "test@example.com", "password": "password123"}`,
			mockSetup: func(m *MockUserService) {
				m.On("LoginUser", mock.Anything, "test@example.com", "password123").
					Return("login-token", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"token":"login-token"}`,
		},
		{
			name:        "invalid credentials",
			requestBody: `{"email": "test@example.com", "password": "wrong"}`,
			mockSetup: func(m *MockUserService) {
				m.On("LoginUser", mock.Anything, "test@example.com", "wrong").
					Return("", errors.New("invalid credentials"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid credentials"}`,
		},
		{
			name:        "invalid JSON",
			requestBody: `invalid json`,
			mockSetup: func(m *MockUserService) {
				// No mock setup needed
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}

			handler := RegisterUserHandlers(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/login", strings.NewReader(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.PostLogin(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandlers_PostRegister(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*MockUserService)
		expectedCode int
		expectedBody string
	}{
		{
			name:        "successful registration",
			requestBody: `{"email": "new@example.com", "password": "password123", "role": "moderator"}`,
			mockSetup: func(m *MockUserService) {
				m.On("RegisterUser", mock.Anything, "new@example.com", "password123", "moderator").
					Return("register-token", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"token":"register-token"}`,
		},
		{
			name:        "registration failed",
			requestBody: `{"email": "exists@example.com", "password": "password123", "role": "employee"}`,
			mockSetup: func(m *MockUserService) {
				m.On("RegisterUser", mock.Anything, "exists@example.com", "password123", "employee").
					Return("", errors.New("email already exists"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"email already exists"}`,
		},
		{
			name:        "missing required fields",
			requestBody: `{"email": "incomplete@example.com"}`,
			mockSetup: func(m *MockUserService) {
				m.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return("", errors.New("invalid params"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid params"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}

			handler := RegisterUserHandlers(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/register", strings.NewReader(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.PostRegister(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
