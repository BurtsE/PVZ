package receptions

import (
	"context"
	domain "pvz/internal/domain/receptions"
	"time"
)

func (s *ReceptionService) CreatePoint(ctx context.Context, point domain.PickUpPoint) error {
	return s.receptionRepo.CreatePoint(ctx, point)
}

func (s *ReceptionService) GetPointList(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]*domain.PickUpPoint, error) {
	//TODO implement me
	panic("implement me")
}
