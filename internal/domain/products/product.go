package products

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const (
	Electronics = 1 << iota
	Clothes
	Shoes
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid products")
)

type Product struct {
	id          uuid.UUID
	name        string
	productType uint32
	price       float64
}

func NewProduct(id uuid.UUID, name string, productType string, price float64) (*Product, error) {
	var (
		t   uint32
		err error
	)
	if err = validateProductName(name); err != nil {
		return nil, err
	}
	if err = validateProductPrice(price); err != nil {
		return nil, err
	}
	if t, err = validateProductType(productType); err != nil {
		return nil, err
	}

	return &Product{
		id:          id,
		name:        name,
		productType: t,
		price:       price,
	}, nil
}

func CreateProduct(name string, productType string, price float64) (*Product, error) {
	return NewProduct(uuid.New(), name, productType, price)
}

func (p *Product) ID() uuid.UUID {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Price() float64 {
	return p.price
}

func validateProductName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidProduct)
	}
	return nil
}

func validateProductPrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("%w: price must be greater than 0", ErrInvalidProduct)
	}
	return nil
}

func validateProductType(productType string) (uint32, error) {
	switch productType {
	case "электроника":
		return Electronics, nil
	case "одежда":
		return Clothes, nil
	case "обувь":
		return Shoes, nil
	}
	return 0, fmt.Errorf("%w: products type not supported", ErrInvalidProduct)
}
