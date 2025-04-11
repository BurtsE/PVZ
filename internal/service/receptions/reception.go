package receptions

import (
	"context"
	"github.com/google/uuid"
	domain "pvz/internal/domain/receptions"
)

func (r *ReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReceptionService) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (domain.PickUpPoint, error) {
	//TODO implement me
	panic("implement me")
}
