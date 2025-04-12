package receptions

import (
	"context"
	"github.com/google/uuid"
	app "pvz/internal/application/receptions"
	domain "pvz/internal/domain/receptions"
	"time"
)

var _ app.ReceptionService = (*ReceptionService)(nil)

type ReceptionRepository interface {
	CreatePoint(ctx context.Context, point domain.PickUpPoint) error
	GetPointList(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*domain.PickUpPoint, error)

	CreateReception(ctx context.Context, reception domain.Reception, pvzID uuid.UUID) error
	CloseReception(ctx context.Context, pvzId uuid.UUID) (domain.Reception, error)

	AddProduct(ctx context.Context, product domain.Product, pvzID uuid.UUID) (uuid.UUID, error)
	DeleteLastAddedProduct(ctx context.Context, pvzID uuid.UUID) error
}

type ReceptionService struct {
	receptionRepo ReceptionRepository
}

func NewReceptionService(rr ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		receptionRepo: rr,
	}
}
