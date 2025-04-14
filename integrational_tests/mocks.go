package integrational_tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	domain "pvz/internal/domain/receptions"
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
