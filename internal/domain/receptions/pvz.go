package receptions

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	MOSCOW = "Москва"
	SPB    = "Санкт-Петербург"
	KAZAN  = "Казань"
)

var (
	ErrPickUpPointNotFound = errors.New("point not found")
	ErrInvalidPickUpPoint  = errors.New("invalid pick-up point")
)

type PickUpPoint struct {
	id               uuid.UUID
	registrationDate time.Time
	city             string
	receptions       []Reception
}

func NewPickUpPoint(id uuid.UUID, registrationDate time.Time, city string, receptions []Reception) (*PickUpPoint, error) {

	if err := validatePointCity(city); err != nil {
		return nil, err
	}
	if err := validatePointRegistry(registrationDate); err != nil {
		return nil, err
	}

	return &PickUpPoint{
		id:               id,
		registrationDate: registrationDate,
		city:             city,
		receptions:       receptions,
	}, nil
}

func CreatePickUpPoint(registrationDate time.Time, city string, receptions []Reception) (*PickUpPoint, error) {
	return NewPickUpPoint(uuid.New(), registrationDate, city, receptions)
}

func (p PickUpPoint) ID() uuid.UUID {
	return p.id
}

func (p PickUpPoint) RegistrationDate() time.Time {
	return p.registrationDate
}

func (p PickUpPoint) City() string {
	return p.city
}

func (p PickUpPoint) Receptions() []Reception {
	return p.receptions
}

func validatePointRegistry(registrationDate time.Time) error {
	if registrationDate.Before(time.Now().AddDate(-1, 0, 0)) {
		return fmt.Errorf("%w: registration expired", ErrInvalidPickUpPoint)
	} else if registrationDate.After(time.Now()) {
		return fmt.Errorf("%w: invalid registry date", ErrInvalidPickUpPoint)
	}
	return nil
}
func validatePointCity(city string) error {
	switch city {
	case MOSCOW:
		return nil
	case SPB:
		return nil
	case KAZAN:
		return nil
	}
	return fmt.Errorf("%w: city not supported", ErrInvalidPickUpPoint)
}
