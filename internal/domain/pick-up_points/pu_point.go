package pick_up_points

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	Moscow = 1 << iota
	Spb
	Kazan
)

var (
	ErrPickUpPointNotFound = errors.New("products not found")
	ErrInvalidPickUpPoint  = errors.New("invalid pick-up point")
)

type PickUpPoint struct {
	id               uuid.UUID
	registrationDate time.Time
	city             uint64
}

func NewPickUpPoint(id uuid.UUID, registrationDate time.Time, cityStr string) (*PickUpPoint, error) {
	var (
		city uint64
		err  error
	)
	if city, err = validatePointCity(cityStr); err != nil {
		return nil, err
	}
	if err = validatePointRegistry(registrationDate); err != nil {
		return nil, err
	}

	return &PickUpPoint{
		id:               id,
		registrationDate: registrationDate,
		city:             city,
	}, nil
}

func CreatePickUpPoint(registrationDate time.Time, city string) (*PickUpPoint, error) {
	return NewPickUpPoint(uuid.New(), registrationDate, city)
}

func validatePointRegistry(registrationDate time.Time) error {
	if registrationDate.Before(time.Now().AddDate(-1, 0, 0)) {
		return fmt.Errorf("%w: registration expired", ErrInvalidPickUpPoint)
	} else if registrationDate.After(time.Now()) {
		return fmt.Errorf("%w: invalid registry date", ErrInvalidPickUpPoint)
	}
	return nil
}
func validatePointCity(city string) (uint64, error) {
	switch city {
	case "Москва":
		return Moscow, nil
	case "Санкт-Петербург":
		return Spb, nil
	case "Казань":
		return Kazan, nil
	}
	return 0, fmt.Errorf("%w: city not supported", ErrInvalidPickUpPoint)

}
