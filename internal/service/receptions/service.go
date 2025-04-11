package receptions

import (
	"context"
	app "pvz/internal/application/receptions"
	domain "pvz/internal/domain/receptions"
)

var _ app.ReceptionService = (*ReceptionService)(nil)

type ReceptionRepository interface {
	CreatePoint(ctx context.Context, point domain.PickUpPoint) error
}

type ReceptionService struct {
	receptionRepo ReceptionRepository
}

func NewReceptionService(rr ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		receptionRepo: rr,
	}
}
