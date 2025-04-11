package receptions

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	Electronics = "электроника"
	Clothes     = "одежда"
	Shoes       = "обувь"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid products")
)

type Product struct {
	id          uuid.UUID
	arrivalTime time.Time
	productType string
}

func NewProduct(id uuid.UUID, time time.Time, productType string) (*Product, error) {
	if err := validateProductType(productType); err != nil {
		return nil, err
	}

	return &Product{
		id:          id,
		arrivalTime: time,
		productType: productType,
	}, nil
}

func CreateProduct(productType string) (*Product, error) {
	return NewProduct(uuid.New(), time.Now(), productType)
}

func (p *Product) ID() uuid.UUID {
	return p.id
}

func (p *Product) ArrivalTime() time.Time {
	return p.arrivalTime
}

func (p *Product) ProductType() string {
	return p.productType
}

func validateProductType(productType string) error {
	switch productType {
	case Electronics:
		return nil
	case Clothes:
		return nil
	case Shoes:
		return nil
	}
	return fmt.Errorf("%w: products type not supported", ErrInvalidProduct)
}
