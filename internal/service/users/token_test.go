package users

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "pvz/internal/domain/users"
)

func TestToken_Parse(t *testing.T) {
	secretKey := []byte("test-secret-key")
	tokenService := NewToken(secretKey)

	// Создаем тестовые данные
	userID := uuid.New()
	email := "test@example.com"
	role := "employee"
	passwordHash := []byte("hashed-password")

	t.Run("valid token", func(t *testing.T) {
		// Создаем валидный токен
		user, err := domain.NewUser(userID, email, role, passwordHash)
		require.NoError(t, err)
		tokenString, err := tokenService.Create(user)
		require.NoError(t, err)

		// Парсим токен
		parsedUser, err := tokenService.Parse(tokenString)
		require.NoError(t, err)

		// Проверяем данные
		assert.Equal(t, userID, parsedUser.ID())
		assert.Equal(t, email, parsedUser.Email())
		assert.Equal(t, role, parsedUser.Role())
		assert.Equal(t, passwordHash, parsedUser.PasswordHash())
	})

	t.Run("invalid token - wrong signing method", func(t *testing.T) {
		// Создаем токен с неправильным методом подписи
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ID":            userID.String(),
			"Email":         email,
			"Role":          role,
			"Password_hash": string(passwordHash),
		})
		tokenString, err := token.SignedString([]byte("wrong-key"))
		require.NoError(t, err)
		_, err = tokenService.Parse(tokenString)
		assert.ErrorIs(t, err, jwt.ErrSignatureInvalid)
	})

	t.Run("invalid token - missing claims", func(t *testing.T) {
		// Создаем токен с отсутствующими claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ID":    userID.String(),
			"Email": email,
			// Нет Role и Password_hash
		})
		tokenString, err := token.SignedString(secretKey)
		require.NoError(t, err)

		_, err = tokenService.Parse(tokenString)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("invalid token - expired", func(t *testing.T) {
		// Создаем токен с истекшим сроком
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ID":            userID.String(),
			"Email":         email,
			"Role":          role,
			"Password_hash": string(passwordHash),
			"exp":           time.Now().Add(-time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(secretKey)
		require.NoError(t, err)

		_, err = tokenService.Parse(tokenString)
		assert.Error(t, err)
	})

	t.Run("invalid token - wrong UUID format", func(t *testing.T) {
		// Создаем токен с невалидным UUID
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ID":            "not-a-uuid",
			"Email":         email,
			"Role":          role,
			"Password_hash": string(passwordHash),
		})
		tokenString, err := token.SignedString(secretKey)
		require.NoError(t, err)

		_, err = tokenService.Parse(tokenString)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})
}

func TestToken_Create(t *testing.T) {
	secretKey := []byte("test-secret-key")
	tokenService := NewToken(secretKey)

	userID := uuid.New()
	email := "test@example.com"
	role := "employee"
	passwordHash := []byte("hashed-password")

	t.Run("create valid token", func(t *testing.T) {
		user, err := domain.NewUser(userID, email, role, passwordHash)
		require.NoError(t, err)
		tokenString, err := tokenService.Create(user)
		require.NoError(t, err)

		// Проверяем, что токен может быть распарсен
		parsedUser, err := tokenService.Parse(tokenString)
		require.NoError(t, err)

		assert.Equal(t, userID, parsedUser.ID())
		assert.Equal(t, email, parsedUser.Email())
		assert.Equal(t, role, parsedUser.Role())
		assert.Equal(t, passwordHash, parsedUser.PasswordHash())
	})

	t.Run("create token with empty user", func(t *testing.T) {
		emptyUser, _ := domain.NewUser(uuid.Nil, "", "", nil)
		_, err := tokenService.Create(emptyUser)
		require.NoError(t, err)
		// Проверяем, что токен создается даже с пустыми значениями
	})
}
