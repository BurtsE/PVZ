package receptions

import (
	"context"
	domain "pvz/internal/domain/receptions"
	"time"
)

func (r *ReceptionService) CreatePoint(ctx context.Context, point domain.PickUpPoint) error {
	return r.receptionRepo.CreatePoint(ctx, point)
}

func (r *ReceptionService) GetPointList(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]*domain.PickUpPoint, error) {
	//TODO implement me
	panic("implement me")
}
