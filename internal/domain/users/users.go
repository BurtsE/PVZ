package users

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const (
	Moderator = "moderator"
	Employee  = "employee"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidUser      = errors.New("invalid user")
	ErrUserValidation   = errors.New("validation error")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type User struct {
	id           uuid.UUID
	email        string
	role         string
	passwordHash []byte
}

func NewUser(id uuid.UUID, email, role string, passwordHash []byte) (User, error) {

	if err := validateEmail(email); err != nil {
		return User{}, err
	}
	if err := validateRole(role); err != nil {
		return User{}, err
	}

	return User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
	}, nil
}

func CreateUser(email, role string, passwordHash []byte) (User, error) {
	return NewUser(uuid.New(), email, role, passwordHash)
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) PasswordHash() []byte {
	return u.passwordHash
}

func (u *User) Role() string {
	switch u.role {
	case Employee:
		return "employee"
	case Moderator:
		return "moderator"
	default:
		return ""
	}
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

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("%w: email is required", ErrUserValidation)
	}
	return nil
}

func validateRole(role string) error {
	switch role {
	case Employee:
		return nil
	case Moderator:
		return nil
	}
	return fmt.Errorf("%w: role not supported", ErrInvalidUser)
}
