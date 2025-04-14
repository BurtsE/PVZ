package receptions

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"pvz/generated/openapi"
	domain "pvz/internal/domain/receptions"
	"pvz/internal/domain/users"
	"testing"
	"time"
)

type MockReceptionService struct {
	mock.Mock
}

func (m *MockReceptionService) CreatePoint(ctx context.Context, point domain.PickUpPoint) error {
	args := m.Called(ctx, point)
	return args.Error(0)
}

func (m *MockReceptionService) GetPointList(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]*domain.PickUpPoint, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	return args.Get(0).([]*domain.PickUpPoint), args.Error(1)
}

func (m *MockReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(domain.Reception), args.Error(1)
}

func (m *MockReceptionService) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(domain.Reception), args.Error(1)
}

func (m *MockReceptionService) AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (domain.Product, uuid.UUID, error) {
	args := m.Called(ctx, pvzID, productType)
	return args.Get(0).(domain.Product), args.Get(1).(uuid.UUID), args.Error(2)
}

func (m *MockReceptionService) DeleteLastProductFromReception(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func TestReceptionHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createTestUser := func(role string) users.User {
		user, err := users.NewUser(uuid.New(), "test@example.com", role, []byte("hash"))
		require.NoError(t, err)
		return user
	}

	t.Run("PostProducts", func(t *testing.T) {
		tests := []struct {
			name          string
			userRole      string
			requestBody   string
			mockSetup     func(*MockReceptionService)
			expectedCode  int
			expectedError string
		}{
			{
				name:        "success - employee can add products",
				userRole:    users.Employee,
				requestBody: `{"pvzId": "550e8400-e29b-41d4-a716-446655440000", "type": "electronics"}`,
				mockSetup: func(m *MockReceptionService) {
					pvzID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
					product := domain.Product{}
					receptionID := uuid.New()
					m.On("AddProduct", mock.Anything, pvzID, "electronics").Return(product, receptionID, nil)
				},
				expectedCode: http.StatusOK,
			},
			{
				name:          "unauthorized - moderator cannot add products",
				userRole:      users.Moderator,
				expectedCode:  http.StatusForbidden,
				expectedError: "only employees can add products",
			},
			{
				name:          "invalid pvzId format",
				userRole:      users.Employee,
				requestBody:   `{"pvzId": "invalid", "type": "electronics"}`,
				expectedCode:  http.StatusBadRequest,
				expectedError: "invalid UUID",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.userRole != "" {
					user := createTestUser(tt.userRole)
					c.Set("user", user)
				}

				mockService := new(MockReceptionService)
				if tt.mockSetup != nil {
					tt.mockSetup(mockService)
				}

				handler := RegisterReceptionHandlers(mockService)

				req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(tt.requestBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req

				handler.PostProducts(c)

				assert.Equal(t, tt.expectedCode, w.Code)
				if tt.expectedError != "" {
					assert.Contains(t, w.Body.String(), tt.expectedError)
				}
				mockService.AssertExpectations(t)
			})
		}
	})

	t.Run("GetPvz", func(t *testing.T) {
		now := time.Now()
		oneMonthAgo := now.AddDate(0, -1, 0)

		tests := []struct {
			name          string
			userRole      string
			page          int
			limit         int
			startDate     *time.Time
			endDate       *time.Time
			mockSetup     func(*MockReceptionService)
			expectedCode  int
			expectedError string
		}{
			{
				name:      "success - employee can access",
				userRole:  users.Employee,
				page:      1,
				limit:     10,
				startDate: &oneMonthAgo,
				endDate:   &now,
				mockSetup: func(m *MockReceptionService) {
					m.On("GetPointList", mock.Anything, oneMonthAgo, now, 1, 10).
						Return([]*domain.PickUpPoint{}, nil)
				},
				expectedCode: http.StatusOK,
			},
			{
				name:      "success - moderator can access",
				userRole:  users.Moderator,
				page:      1,
				limit:     10,
				startDate: &oneMonthAgo,
				endDate:   &now,
				mockSetup: func(m *MockReceptionService) {
					m.On("GetPointList", mock.Anything, oneMonthAgo, now, 1, 10).
						Return([]*domain.PickUpPoint{}, nil)
				},
				expectedCode: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.userRole != "" {
					user := createTestUser(tt.userRole)
					c.Set("user", user)
				}

				mockService := new(MockReceptionService)
				if tt.mockSetup != nil {
					tt.mockSetup(mockService)
				}

				handler := RegisterReceptionHandlers(mockService)

				params := openapi.GetPvzParams{
					Page:      &tt.page,
					Limit:     &tt.limit,
					StartDate: tt.startDate,
					EndDate:   tt.endDate,
				}

				handler.GetPvz(c, params)

				assert.Equal(t, tt.expectedCode, w.Code)
				if tt.expectedError != "" {
					assert.Contains(t, w.Body.String(), tt.expectedError)
				}
				mockService.AssertExpectations(t)
			})
		}
	})

	t.Run("PostPvz", func(t *testing.T) {
		tests := []struct {
			name          string
			userRole      string
			requestBody   string
			mockSetup     func(*MockReceptionService)
			expectedCode  int
			expectedError string
		}{
			{
				name:        "success - moderator can add points",
				userRole:    users.Moderator,
				requestBody: `{"id": "550e8400-e29b-41d4-a716-446655440000", "registrationDate": "2025-01-01T00:00:00Z", "city": "Москва"}`,
				mockSetup: func(m *MockReceptionService) {
					m.On("CreatePoint", mock.Anything, mock.AnythingOfType("receptions.PickUpPoint")).
						Return(nil)
				},
				expectedCode: http.StatusOK,
			},
			{
				name:          "unauthorized - employee cannot add points",
				userRole:      users.Employee,
				expectedCode:  http.StatusForbidden,
				expectedError: "only moderator can add points",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.userRole != "" {
					user := createTestUser(tt.userRole)
					c.Set("user", user)
				}

				mockService := new(MockReceptionService)
				if tt.mockSetup != nil {
					tt.mockSetup(mockService)
				}

				handler := RegisterReceptionHandlers(mockService)

				req, _ := http.NewRequest("POST", "/pvz", bytes.NewBufferString(tt.requestBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req

				handler.PostPvz(c)

				assert.Equal(t, tt.expectedCode, w.Code)
				if tt.expectedError != "" {
					assert.Contains(t, w.Body.String(), tt.expectedError)
				}
				mockService.AssertExpectations(t)
			})
		}
	})
	t.Run("PostPvzPvzIdCloseLastReception", func(t *testing.T) {
		tests := []struct {
			name          string
			userRole      string
			mockSetup     func(*MockReceptionService)
			expectedCode  int
			expectedError string
		}{
			{
				name:     "success - employee can close reception",
				userRole: users.Employee,
				mockSetup: func(m *MockReceptionService) {
					m.On("CloseLastReception", mock.Anything, mock.AnythingOfType("uuid.UUID")).
						Return(domain.Reception{}, nil)
				},
				expectedCode: http.StatusOK,
			},
			{
				name:          "unauthorized - employee cannot close reception",
				userRole:      "",
				expectedCode:  http.StatusForbidden,
				expectedError: "only employees can close receptions",
			},
		}
		pvzID, err := uuid.Parse("3fa85f64-5717-4562-b3fc-2c963f66afa7")
		require.NoError(t, err)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.userRole != "" {
					user := createTestUser(tt.userRole)
					c.Set("user", user)
				}

				mockService := new(MockReceptionService)
				if tt.mockSetup != nil {
					tt.mockSetup(mockService)
				}

				handler := RegisterReceptionHandlers(mockService)

				req, _ := http.NewRequest("POST", fmt.Sprintf("/pvz/%s/close_last_reception", pvzID.String()), nil)
				c.Request = req

				handler.PostPvzPvzIdCloseLastReception(c, pvzID)

				assert.Equal(t, tt.expectedCode, w.Code)
				if tt.expectedError != "" {
					assert.Contains(t, w.Body.String(), tt.expectedError)
				}
				mockService.AssertExpectations(t)
			})
		}
	})
	t.Run("PostPvzPvzIdDeleteLastProduct", func(t *testing.T) {
		tests := []struct {
			name          string
			userRole      string
			mockSetup     func(*MockReceptionService)
			expectedCode  int
			expectedError string
		}{
			{
				name:     "success - employee can delete product",
				userRole: users.Employee,
				mockSetup: func(m *MockReceptionService) {
					m.On("DeleteLastProductFromReception", mock.Anything, mock.AnythingOfType("uuid.UUID")).
						Return(nil)
				},
				expectedCode: http.StatusOK,
			},
			{
				name:          "unauthorized - employee cannot delete product",
				userRole:      "",
				expectedCode:  http.StatusForbidden,
				expectedError: "only employees can delete products",
			},
		}
		pvzID, err := uuid.Parse("3fa85f64-5717-4562-b3fc-2c963f66afa7")
		require.NoError(t, err)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.userRole != "" {
					user := createTestUser(tt.userRole)
					c.Set("user", user)
				}

				mockService := new(MockReceptionService)
				if tt.mockSetup != nil {
					tt.mockSetup(mockService)
				}

				handler := RegisterReceptionHandlers(mockService)

				req, _ := http.NewRequest("POST", fmt.Sprintf("/pvz/%s/delete_last_product", pvzID.String()), nil)
				c.Request = req

				handler.PostPvzPvzIdDeleteLastProduct(c, pvzID)

				assert.Equal(t, tt.expectedCode, w.Code)
				if tt.expectedError != "" {
					assert.Contains(t, w.Body.String(), tt.expectedError)
				}
				mockService.AssertExpectations(t)
			})
		}
	})
	t.Run("PostReceptions", func(t *testing.T) {
		tests := []struct {
			name          string
			userRole      string
			requestBody   string
			mockSetup     func(*MockReceptionService)
			expectedCode  int
			expectedError string
		}{
			{
				name:        "success - employee can delete product",
				userRole:    users.Employee,
				requestBody: `{"pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa7"}`,
				mockSetup: func(m *MockReceptionService) {
					m.On("CreateReception", mock.Anything, mock.AnythingOfType("uuid.UUID")).
						Return(domain.Reception{}, nil)
				},
				expectedCode: http.StatusOK,
			},
			{
				name:          "unauthorized - employee cannot add points",
				userRole:      "",
				expectedCode:  http.StatusForbidden,
				expectedError: "only employees can open receptions",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Setup similar to previous tests
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.userRole != "" {
					user := createTestUser(tt.userRole)
					c.Set("user", user)
				}

				mockService := new(MockReceptionService)
				if tt.mockSetup != nil {
					tt.mockSetup(mockService)
				}

				handler := RegisterReceptionHandlers(mockService)

				req, _ := http.NewRequest("POST", "/receptions", bytes.NewBufferString(tt.requestBody))
				c.Request = req

				handler.PostReceptions(c)

				assert.Equal(t, tt.expectedCode, w.Code)
				if tt.expectedError != "" {
					assert.Contains(t, w.Body.String(), tt.expectedError)
				}
				mockService.AssertExpectations(t)
			})
		}
	})
}

func TestConvertProductToResponse(t *testing.T) {
	now := time.Date(2025, time.May, 15, 12, 0, 0, 0, time.UTC)
	productID := uuid.New()
	receptionID := uuid.New().String()
	product, err := domain.NewProduct(productID, now, domain.Electronics)
	require.NoError(t, err)

	expected := gin.H{
		"id":          productID,
		"dateTime":    now,
		"type":        domain.Electronics,
		"receptionId": receptionID,
	}

	result := convertProductToResponse(*product, receptionID)
	assert.Equal(t, expected["id"], result["id"])
	assert.Equal(t, expected["type"], result["type"])
	assert.Equal(t, expected["receptionId"], result["receptionId"])
	resultTime, ok := result["dateTime"].(time.Time)
	assert.True(t, ok)
	assert.WithinDuration(t, now, resultTime, time.Millisecond)
}

func TestConvertReceptionToResponse_WithProducts(t *testing.T) {
	now := time.Now()
	receptionID := uuid.New()
	pvzID := uuid.New().String()
	products := []domain.Product{}
	for i := range 2 {
		diff := time.Duration(i) * time.Hour
		product, err := domain.NewProduct(uuid.New(), now.Add(diff), domain.Clothes)
		require.NoError(t, err)
		products = append(products, *product)
	}
	reception, err := domain.NewReception(receptionID, domain.IN_PROGRESS, now, products)
	require.NoError(t, err)

	result := convertReceptionToResponse(*reception, pvzID)

	assert.Equal(t, receptionID, result["id"])
	assert.Equal(t, now, result["dateTime"])
	assert.Equal(t, pvzID, result["pvzId"])
	assert.Equal(t, domain.IN_PROGRESS, result["status"])

	productsResult, ok := result["products"].([]gin.H)
	assert.True(t, ok)
	assert.Len(t, productsResult, 2)

	for i, p := range productsResult {
		assert.Equal(t, products[i].ID(), p["id"])
		assert.Equal(t, products[i].ArrivalTime(), p["dateTime"])
		assert.Equal(t, products[i].ProductType(), p["type"])
		assert.Equal(t, receptionID.String(), p["receptionId"])
	}
}

func TestConvertReceptionToResponse_WithoutProducts(t *testing.T) {
	now := time.Now()
	receptionID := uuid.New()
	pvzID := uuid.New().String()
	reception, err := domain.NewReception(receptionID, domain.IN_PROGRESS, now, []domain.Product{})
	require.NoError(t, err)

	result := convertReceptionToResponse(*reception, pvzID)

	assert.Equal(t, receptionID, result["id"])
	assert.Equal(t, now, result["dateTime"])
	assert.Equal(t, pvzID, result["pvzId"])
	assert.Equal(t, domain.IN_PROGRESS, result["status"])
	assert.Nil(t, result["products"])
}

func TestConvertPointListToResponse(t *testing.T) {
	now := time.Now()
	pickUpPointID := uuid.New()
	city := domain.KAZAN

	products := []domain.Product{}
	for i := range 3 {
		diff := time.Duration(i) * time.Hour
		product, err := domain.NewProduct(uuid.New(), now.Add(diff), domain.Shoes)
		require.NoError(t, err)
		products = append(products, *product)
	}
	reception1, err := domain.NewReception(uuid.New(), domain.IN_PROGRESS, now.Add(-3*time.Hour), []domain.Product{products[0]})
	require.NoError(t, err)
	reception2, err := domain.NewReception(uuid.New(), domain.CLOSED, now.Add(-4*time.Hour), []domain.Product{products[1]})
	require.NoError(t, err)

	pickUpPoint, err := domain.NewPickUpPoint(pickUpPointID, now.AddDate(0, 0, -1), domain.KAZAN, []domain.Reception{*reception1, *reception2})
	require.NoError(t, err)

	points := []*domain.PickUpPoint{pickUpPoint}

	result := convertPointListToResponse(points)

	assert.Len(t, result, 1)
	pointResult := result[0]

	pvz, ok := pointResult["pvz"].(gin.H)
	assert.True(t, ok)
	assert.Equal(t, pickUpPointID, pvz["id"])
	assert.Equal(t, city, pvz["city"])
	assert.Equal(t, now.AddDate(0, 0, -1), pvz["registrationDate"])

	receptions, ok := pointResult["receptions"].([]gin.H)
	assert.True(t, ok)
	assert.Len(t, receptions, 2)

	for i, r := range receptions {
		var expectedReception domain.Reception
		if i == 0 {
			expectedReception = *reception1
		} else {
			expectedReception = *reception2
		}

		assert.Equal(t, expectedReception.ID(), r["id"])
		assert.Equal(t, expectedReception.InitTime(), r["dateTime"])
		assert.Equal(t, pickUpPointID.String(), r["pvzId"])
		assert.Equal(t, expectedReception.Status(), r["status"])

		if len(expectedReception.Products()) > 0 {
			products, ok := r["products"].([]gin.H)
			assert.True(t, ok)
			assert.Len(t, products, 1)
			assert.Equal(t, expectedReception.Products()[0].ID(), products[0]["id"])
		}
	}
}

func TestConvertPointListToResponse_Empty(t *testing.T) {
	result := convertPointListToResponse([]*domain.PickUpPoint{})
	assert.Empty(t, result)
}
