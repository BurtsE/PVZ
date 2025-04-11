package receptions

import (
	"context"
	"github.com/google/uuid"
	domain "pvz/internal/domain/receptions"
)

func (s *ReceptionService) AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ReceptionService) DeleteLastProductFromReception(ctx context.Context, receptionID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
