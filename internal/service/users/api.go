package users

import (
	"context"
	domain "pvz/internal/domain/users"
)

func (s *UserService) RegisterUser(ctx context.Context, email, password, role string) (string, error) {
	err := checkPassword(password)
	if err != nil {
		return "", err
	}
	hash, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	user, err := domain.CreateUser(email, role, hash)
	if err != nil {
		return "", err
	}
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

func (s *UserService) ValidateUser(_ context.Context, token string) (domain.User, error) {
	user, err := s.token.Parse(token)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetUser(ctx, email)
	if err != nil {
		return "", err
	}
	err = compareHashAndPassword(user.PasswordHash(), password)
	if err != nil {
		return "", err
	}
	token, err := s.token.Create(user)
	if err != nil {
		return "", err
	}
	return token, nil

}
