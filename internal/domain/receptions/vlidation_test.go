package receptions

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateReceptionStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected error
	}{
		{
			name:     "valid in_progress status",
			status:   IN_PROGRESS,
			expected: nil,
		},
		{
			name:     "valid closed status",
			status:   CLOSED,
			expected: nil,
		},
		{
			name:     "invalid status",
			status:   "invalid_status",
			expected: errors.New("invalid pick-up point: status not supported"),
		},
		{
			name:     "empty status",
			status:   "",
			expected: errors.New("invalid pick-up point: status not supported"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReceptionStatus(tt.status)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func TestValidatePointRegistry(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name             string
		registrationDate time.Time
		expected         error
	}{
		{
			name:             "valid current date",
			registrationDate: now,
			expected:         nil,
		},
		{
			name:             "valid date within 1 year",
			registrationDate: now.AddDate(0, -6, 0), // 6 months ago
			expected:         nil,
		},
		{
			name:             "expired registration (older than 1 year)",
			registrationDate: now.AddDate(-2, 0, 0), // 2 years ago
			expected:         errors.New("invalid pick-up point: registration expired"),
		},
		{
			name:             "future date",
			registrationDate: now.AddDate(0, 0, 1), // tomorrow
			expected:         errors.New("invalid pick-up point: invalid registry date"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePointRegistry(tt.registrationDate)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func TestValidatePointCity(t *testing.T) {
	tests := []struct {
		name     string
		city     string
		expected error
	}{
		{
			name:     "valid Moscow",
			city:     MOSCOW,
			expected: nil,
		},
		{
			name:     "valid Saint Petersburg",
			city:     SPB,
			expected: nil,
		},
		{
			name:     "valid Kazan",
			city:     KAZAN,
			expected: nil,
		},
		{
			name:     "invalid city",
			city:     "Новосибирск",
			expected: errors.New("invalid pick-up point: city not supported"),
		},
		{
			name:     "empty city",
			city:     "",
			expected: errors.New("invalid pick-up point: city not supported"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePointCity(tt.city)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func TestValidateProductType(t *testing.T) {
	tests := []struct {
		name        string
		productType string
		expected    error
	}{
		{
			name:        "valid electronics",
			productType: Electronics,
			expected:    nil,
		},
		{
			name:        "valid clothes",
			productType: Clothes,
			expected:    nil,
		},
		{
			name:        "valid shoes",
			productType: Shoes,
			expected:    nil,
		},
		{
			name:        "invalid product type",
			productType: "мебель",
			expected:    errors.New("invalid product: products type not supported"),
		},
		{
			name:        "empty product type",
			productType: "",
			expected:    errors.New("invalid product: products type not supported"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProductType(tt.productType)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}
