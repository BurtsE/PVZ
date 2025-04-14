package receptions

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	domain "pvz/internal/domain/receptions"
)

type MockReceptionRepo struct {
	mock.Mock
}

func (m *MockReceptionRepo) AddProduct(ctx context.Context, product domain.Product, pvzID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, product, pvzID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockReceptionRepo) DeleteLastAddedProduct(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func (m *MockReceptionRepo) CreatePoint(ctx context.Context, point domain.PickUpPoint) error {
	args := m.Called(ctx, point)
	return args.Error(0)
}

func (m *MockReceptionRepo) GetPointList(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*domain.PickUpPoint, error) {
	args := m.Called(ctx, startDate, endDate, limit, offset)
	return args.Get(0).([]*domain.PickUpPoint), args.Error(1)
}

func (m *MockReceptionRepo) CreateReception(ctx context.Context, reception domain.Reception, pvzID uuid.UUID) error {
	args := m.Called(ctx, reception, pvzID)
	return args.Error(0)
}

func (m *MockReceptionRepo) CloseReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(domain.Reception), args.Error(1)
}

func TestReceptionService_AddProduct(t *testing.T) {
	t.Run("successful product addition", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()
		productID := uuid.New()
		expectedProduct, _ := domain.NewProduct(productID, time.Now(), domain.Electronics)
		expectedReceptionID := uuid.New()

		mockRepo.On("AddProduct", ctx, mock.Anything, pvzID).Return(expectedReceptionID, nil)

		product, receptionID, err := service.AddProduct(ctx, pvzID, domain.Electronics)
		require.NoError(t, err)
		assert.Equal(t, (*expectedProduct).ProductType(), product.ProductType())
		assert.Equal(t, (*expectedProduct).ArrivalTime().Add(500*time.Millisecond).After(product.ArrivalTime()), true)

		assert.Equal(t, expectedReceptionID, receptionID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid product type", func(t *testing.T) {
		service := &ReceptionService{}
		_, _, err := service.AddProduct(context.Background(), uuid.New(), "")
		assert.Error(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()
		mockRepo.On("AddProduct", ctx, mock.Anything, pvzID).Return(uuid.Nil, assert.AnError)

		_, _, err := service.AddProduct(ctx, pvzID, domain.Clothes)
		assert.Error(t, err)
	})
}

func TestReceptionService_DeleteLastProductFromReception(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()
		mockRepo.On("DeleteLastAddedProduct", ctx, pvzID).Return(nil)

		err := service.DeleteLastProductFromReception(ctx, pvzID)
		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()
		mockRepo.On("DeleteLastAddedProduct", ctx, pvzID).Return(assert.AnError)

		err := service.DeleteLastProductFromReception(ctx, pvzID)
		assert.Error(t, err)
	})
}

func TestReceptionService_CreatePoint(t *testing.T) {
	t.Run("successful point creation", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)
		id := uuid.New()
		point, _ := domain.NewPickUpPoint(id, time.Now(), domain.MOSCOW, nil)

		mockRepo.On("CreatePoint", ctx, *point).Return(nil)

		err := service.CreatePoint(ctx, *point)
		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		id := uuid.New()
		point, _ := domain.NewPickUpPoint(id, time.Now(), domain.MOSCOW, nil)

		mockRepo.On("CreatePoint", ctx, *point).Return(assert.AnError)

		err := service.CreatePoint(ctx, *point)
		assert.Error(t, err)
	})
}

func TestReceptionService_GetPointList(t *testing.T) {
	t.Run("successful list retrieval", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		now := time.Now()
		expectedPoints := []*domain.PickUpPoint{}
		for range 2 {
			point, _ := domain.CreatePickUpPoint(time.Now(), "Moscow", nil)
			expectedPoints = append(expectedPoints, point)
		}

		mockRepo.On("GetPointList", ctx, now.AddDate(0, -1, 0), now, 10, 0).
			Return(expectedPoints, nil)

		points, err := service.GetPointList(ctx, now.AddDate(0, -1, 0), now, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, expectedPoints, points)
	})

	t.Run("empty result", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		now := time.Now()
		mockRepo.On("GetPointList", ctx, now.AddDate(0, -1, 0), now, 10, 0).
			Return([]*domain.PickUpPoint{}, nil)

		points, err := service.GetPointList(ctx, now.AddDate(0, -1, 0), now, 1, 10)
		require.NoError(t, err)
		assert.Empty(t, points)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		now := time.Now()
		mockRepo.On("GetPointList", ctx, now.AddDate(0, -1, 0), now, 10, 0).
			Return([]*domain.PickUpPoint{}, assert.AnError)

		_, err := service.GetPointList(ctx, now.AddDate(0, -1, 0), now, 1, 10)
		assert.Error(t, err)
	})
}

func TestReceptionService_CreateReception(t *testing.T) {
	t.Run("successful reception creation", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()

		mockRepo.On("CreateReception", ctx, mock.Anything, pvzID).Return(nil)

		reception, err := service.CreateReception(ctx, pvzID)
		require.NoError(t, err)
		assert.Nil(t, reception.Products())
		assert.Equal(t, domain.IN_PROGRESS, reception.Status())
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()

		mockRepo.On("CreateReception", ctx, mock.Anything, pvzID).Return(assert.AnError)

		_, err := service.CreateReception(ctx, pvzID)
		assert.Error(t, err)
	})
}

func TestReceptionService_CloseLastReception(t *testing.T) {
	t.Run("successful reception closing", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()
		expectedReception, _ := domain.CreateReception("in_progress", time.Now(), nil)

		mockRepo.On("CloseReception", ctx, pvzID).Return(*expectedReception, nil)

		reception, err := service.CloseLastReception(ctx, pvzID)
		require.NoError(t, err)
		assert.Equal(t, *expectedReception, reception)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockReceptionRepo)
		service := NewReceptionService(mockRepo)

		pvzID := uuid.New()
		mockRepo.On("CloseReception", ctx, pvzID).Return(domain.Reception{}, assert.AnError)

		_, err := service.CloseLastReception(ctx, pvzID)
		assert.Error(t, err)
	})
}
