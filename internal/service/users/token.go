package users

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	domain "pvz/internal/domain/users"
)

var ErrInvalidToken = fmt.Errorf("token is invalid")

type Token struct {
	secretKey []byte
}

func NewToken(secretKey []byte) *Token {
	return &Token{secretKey: secretKey}
}
func (t *Token) Parse(tokenString string) (domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return t.secretKey, nil
	})
	if err != nil {
		return domain.User{}, err
	}

	var (
		claims       jwt.MapClaims
		id           uuid.UUID
		email        string
		passwordHash []byte
		role         string
		ok           bool
	)
	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return domain.User{}, ErrInvalidToken
	}
	if id, ok = claims["ID"].(uuid.UUID); !ok {
		return domain.User{}, ErrInvalidToken
	}
	if email, ok = claims["Email"].(string); !ok {
		return domain.User{}, ErrInvalidToken
	}
	if role, ok = claims["Role"].(string); !ok {
		return domain.User{}, ErrInvalidToken
	}
	if passwordHash, ok = claims["Password_hash"].([]byte); !ok {
		return domain.User{}, ErrInvalidToken
	}

	return domain.NewUser(id, email, role, passwordHash)
}

func (t *Token) Create(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":            user.ID(),
		"Email":         user.Email(),
		"Role":          user.Role(),
		"Password_hash": user.PasswordHash(),
	})
	return token.SignedString(t.secretKey)
}
