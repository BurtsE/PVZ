package application

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	userHandlers "pvz/internal/application/users"
	"pvz/internal/domain/users"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(ctx context.Context, email, password, role string) (token string, err error) {
	return "", nil
}

func (m *MockUserService) LoginUser(ctx context.Context, email, password string) (token string, err error) {
	return "", nil
}

func (m *MockUserService) ValidateUser(ctx context.Context, token string) (users.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(users.User), args.Error(1)
}

func TestLoginMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		mockUser       users.User
		mockError      error
		expectUserSet  bool
		expectNextCall bool
	}{
		{
			name:           "no authorization header",
			authHeader:     "",
			expectUserSet:  false,
			expectNextCall: true,
		},
		{
			name:           "malformed authorization header - no space",
			authHeader:     "Bearer",
			expectUserSet:  false,
			expectNextCall: true,
		},
		{
			name:           "malformed authorization header - wrong scheme",
			authHeader:     "Basic abc123",
			expectUserSet:  false,
			expectNextCall: true,
		},
		{
			name:           "valid token but validation fails",
			authHeader:     "Bearer invalid-token",
			mockError:      assert.AnError,
			expectUserSet:  false,
			expectNextCall: true,
		},
		{
			name:           "valid token and validation succeeds",
			authHeader:     "Bearer valid-token",
			mockUser:       users.User{},
			expectUserSet:  true,
			expectNextCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			mockService := new(MockUserService)
			server := &httpServer{
				UserHandlers: &userHandlers.UserHandlers{
					Service: mockService,
				},
			}

			var userInContext *users.User
			var nextCalled bool
			testHandler := func(c *gin.Context) {
				nextCalled = true
				if u, exists := c.Get("user"); exists {
					user := u.(users.User)
					userInContext = &user
				}
				c.String(http.StatusOK, "ok")
			}

			r.Use(server.LoginMiddleware())
			r.GET("/test", testHandler)

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			if strings.HasPrefix(tt.authHeader, "Bearer ") {
				token := strings.TrimPrefix(tt.authHeader, "Bearer ")
				mockService.On("ValidateUser", mock.Anything, token).
					Return(tt.mockUser, tt.mockError)
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.expectNextCall, nextCalled, "Next handler call mismatch")

			if tt.expectUserSet {
				require.NotNil(t, userInContext, "Expected user to be set in context")
				assert.Equal(t, tt.mockUser, *userInContext)
			} else {
				assert.Nil(t, userInContext, "Expected no user in context")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestLoginMiddleware_DeferNext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	mockService := new(MockUserService)
	server := &httpServer{
		UserHandlers: &userHandlers.UserHandlers{
			Service: mockService,
		},
	}

	executionOrder := []string{}

	r.Use(func(c *gin.Context) {
		executionOrder = append(executionOrder, "middleware-start")
		c.Next()
		executionOrder = append(executionOrder, "middleware-end")
	})

	r.Use(server.LoginMiddleware())

	r.Use(func(c *gin.Context) {
		executionOrder = append(executionOrder, "second-middleware")
		c.Next()
	})

	r.GET("/test", func(c *gin.Context) {
		executionOrder = append(executionOrder, "handler")
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	expectedOrder := []string{
		"middleware-start",
		"second-middleware",
		"handler",
		"middleware-end",
	}
	assert.Equal(t, expectedOrder, executionOrder, "Middleware execution order incorrect")
}
