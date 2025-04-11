package users

import (
	"context"
	"github.com/google/uuid"
	userService "pvz/internal/application/users"
	"pvz/internal/domain/users"
)

var _ userService.UserService = (*UserService)(nil)

type UserRepository interface {
	GetUser(ctx context.Context, email string) (users.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (users.User, error)
	CreateUser(ctx context.Context, user users.User) error
}

type UserService struct {
	userRepo UserRepository
	token    *Token
}

func NewUserService(ur UserRepository, token *Token) *UserService {
	return &UserService{
		userRepo: ur,
		token:    token,
	}
}
