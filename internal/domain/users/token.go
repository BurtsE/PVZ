package users

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = fmt.Errorf("token is invalid")

type Token struct {
	secretKey []byte
}

func NewToken(secretKey []byte) *Token {
	return &Token{secretKey: secretKey}
}
func (t *Token) Parse(tokenString string) (User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return t.secretKey, nil
	})
	if err != nil {
		return User{}, err
	}

	var (
		claims jwt.MapClaims
		user   User
		ok     bool
	)
	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return User{}, ErrInvalidToken
	}
	if user.id, ok = claims["ID"].(uuid.UUID); !ok {
		return User{}, ErrInvalidToken
	}
	if user.name, ok = claims["Username"].(string); !ok {
		return User{}, ErrInvalidToken
	}

	return user, nil
}

func (t *Token) Create(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   user.id,
		"Role": user.role,
	})
	return token.SignedString(t.secretKey)
}
