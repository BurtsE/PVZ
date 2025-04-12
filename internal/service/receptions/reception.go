package receptions

import (
	"context"
	"github.com/google/uuid"
	domain "pvz/internal/domain/receptions"
	"time"
)

func (s *ReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error) {
	reception, err := domain.CreateReception("in_progress", time.Now(), nil)
	if err != nil {
		return domain.Reception{}, err
	}
	err = s.receptionRepo.CreateReception(ctx, *reception, pvzID)
	if err != nil {
		return domain.Reception{}, err
	}
	return *reception, err
}

func (s *ReceptionService) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error) {
	return s.receptionRepo.CloseReception(ctx, pvzID)
}
