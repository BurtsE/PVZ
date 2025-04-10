package users

import (
	"context"
	"pvz/internal/domain/users"
)

func (s *UserService) RegisterUser(ctx context.Context, user users.User) (string, error) {
	token, err := s.token.Create(user)
	if err != nil {
		return "", err
	}
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) ValidateUser(_ context.Context, token string) (users.User, error) {
	user, err := s.token.Parse(token)
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}
