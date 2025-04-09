package receptions

import (
	"context"
	"github.com/google/uuid"
	"pvz/internal/domain/users"
)

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
}
