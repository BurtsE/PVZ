package users

import (
	"context"
	"github.com/google/uuid"
	"pvz/internal/application"
	"pvz/internal/domain/users"
)

var _ application.UserService = (*UserService)(nil)

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (users.User, error)
	CreateUser(ctx context.Context, user users.User) error
}

type UserService struct {
	userRepo UserRepository
	token    users.Token
}

func NewReceptionService(ur UserRepository, token users.Token) *UserService {
	return &UserService{
		userRepo: ur,
		token:    token,
	}
}
