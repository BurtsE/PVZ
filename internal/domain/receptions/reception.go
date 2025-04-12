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
	id     uuid.UUID
	status string
	//pvzId    uuid.UUID
	registrationDate time.Time
	products         []Product
}

func NewReception(id uuid.UUID, status string, initTime time.Time, products []Product) (*Reception, error) {
	if err := validateReceptionStatus(status); err != nil {
		return nil, err
	}
	return &Reception{
		id:               id,
		status:           status,
		registrationDate: initTime,
		products:         products,
	}, nil
}

func CreateReception(status string, initTime time.Time, products []Product) (*Reception, error) {
	return NewReception(
		uuid.New(),
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

func (r Reception) InitTime() time.Time {
	return r.registrationDate
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
