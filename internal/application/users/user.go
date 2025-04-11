package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz/generated/openapi"
	"pvz/internal/domain/users"
)

type UserService interface {
	RegisterUser(ctx context.Context, email, password, role string) (token string, err error)
	ValidateUser(ctx context.Context, token string) (users.User, error)
	LoginUser(ctx context.Context, email, password string) (token string, err error)
}

type UserHandlers struct {
	s UserService
}

func RegisterUserHandlers(s UserService) *UserHandlers {
	return &UserHandlers{s}
}

func (u *UserHandlers) PostDummyLogin(c *gin.Context) {
	body := openapi.PostDummyLoginJSONBody{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := u.s.RegisterUser(c, "dummy@example.com", "dummy_password", string(body.Role))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (u *UserHandlers) PostLogin(c *gin.Context) {
	body := openapi.PostLoginJSONBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := u.s.LoginUser(c, string(body.Email), body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
func (u *UserHandlers) PostRegister(c *gin.Context) {
	body := openapi.PostRegisterJSONBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := u.s.RegisterUser(c, string(body.Email), body.Password, string(body.Role))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
