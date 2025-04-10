package users

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	Moderator = 1 << iota
	Employee
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidUser      = errors.New("invalid user")
	ErrUserValidation   = errors.New("validation error")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type User struct {
	id    uuid.UUID
	name  string
	email string
	role  byte
}

func NewUser(id uuid.UUID, name, email, roleStr string) (User, error) {
	var (
		role byte
		err  error
	)
	if err = validateUsername(name); err != nil {
		return User{}, err
	}
	if err = validateEmail(email); err != nil {
		return User{}, err
	}

	if role, err = validateRole(roleStr); err != nil {
		return User{}, err
	}

	return User{
		id:    id,
		name:  name,
		email: email,
		role:  role,
	}, nil
}

func CreateUser(name, email, role string) (User, error) {
	return NewUser(uuid.New(), name, email, role)
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Role() byte {
	return u.role
}

func (u *User) SendToEmail(_ string) error {
	return errors.New("not implemented")
}

func (u *User) ChangeEmail(email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	u.email = email
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("%w: name is required", ErrUserValidation)
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("%w: email is required", ErrUserValidation)
	}
	return nil
}

func validateRole(role string) (byte, error) {
	switch role {
	case "employee":
		return Employee, nil
	case "moderator":
		return Moderator, nil
	}
	return 0, fmt.Errorf("%w: role not supported", ErrInvalidUser)
}
