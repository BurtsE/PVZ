package receptions

import (
	"context"
	"github.com/google/uuid"
	app "pvz/internal/application/receptions"
	domain "pvz/internal/domain/receptions"
)

var _ app.ReceptionService = (*ReceptionService)(nil)

type ReceptionRepository interface {
	CreatePoint(ctx context.Context, point domain.PickUpPoint) error
	CreateReception(ctx context.Context, reception domain.Reception) error
	CloseReception(ctx context.Context, pvzId uuid.UUID) (domain.Reception, error)
}

type ReceptionService struct {
	receptionRepo ReceptionRepository
}

func NewReceptionService(rr ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		receptionRepo: rr,
	}
}
