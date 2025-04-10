package application

import (
	"context"
	"pvz/internal/domain/users"
)

type UserService interface {
	RegisterUser(ctx context.Context, user users.User) (token string, err error)
	ValidateUser(ctx context.Context, token string) (users.User, error)
}
