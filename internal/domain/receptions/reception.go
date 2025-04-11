package receptions

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	IN_PROGRESS = "in_progress"
	CLOSED      = "close"
)

type Reception struct {
	id       uuid.UUID
	status   string
	pvzId    uuid.UUID
	initTime time.Time
	products []Product
}

func NewReception(id, pvzId uuid.UUID, status string, initTime time.Time, products []Product) (*Reception, error) {
	return &Reception{
		id:       id,
		status:   status,
		pvzId:    pvzId,
		initTime: initTime,
		products: products,
	}, nil
}

func CreateReception(pvzId uuid.UUID, status string, initTime time.Time, products []Product) (*Reception, error) {
	return NewReception(
		uuid.New(),
		pvzId,
		status,
		initTime,
		products,
	)
}

func (r Reception) ID() uuid.UUID {
	return r.id
}

func (r Reception) Status() string {
	return r.status
}

func (r Reception) PvzID() uuid.UUID {
	return r.pvzId
}

func (r Reception) InitTime() time.Time {
	return r.initTime
}

func (r Reception) Products() []Product {
	return r.products
}

func validateReceptionStatus(status string) error {
	switch status {
	case IN_PROGRESS:
		return nil
	case CLOSED:
		return nil
	}
	return fmt.Errorf("%w: status not supported", ErrInvalidPickUpPoint)

}
