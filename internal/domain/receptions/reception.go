package receptions

import (
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	id       uuid.UUID
	status   string
	pvzId    uuid.UUID
	initTime time.Time
	products []Product
}

func NewReception(id, pvzId uuid.UUID, status string, initTime time.Time) (*Reception, error) {
	return &Reception{
		id:       id,
		status:   status,
		pvzId:    pvzId,
		initTime: initTime,
	}, nil
}

func CreateReception(pvzId uuid.UUID, status string, initTime time.Time) (*Reception, error) {
	return NewReception(
		uuid.New(),
		pvzId,
		status,
		initTime,
	)
}
