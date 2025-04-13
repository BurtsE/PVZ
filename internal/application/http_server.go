package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz/generated/openapi"
	"pvz/internal/application/receptions"
	"pvz/internal/application/users"
	"pvz/internal/config"
	"strings"
	"time"
)

var _ openapi.ServerInterface = (*httpServer)(nil)

type httpServer struct {
	*users.UserHandlers
	*receptions.ReceptionHandlers
}

func SetupHTTPServer(userService users.UserService, receptionsService receptions.ReceptionService) *http.Server {
	port := config.GetApplicationPort()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	server := &httpServer{
		UserHandlers:      users.RegisterUserHandlers(userService),
		ReceptionHandlers: receptions.RegisterReceptionHandlers(receptionsService),
	}
	r.Use(server.LoginMiddleware())
	//registry := middleware.NewMiddlewareRegistry()
	//registry.Register("PostPvz", server.AdminRequiLoredMiddleware())
	openapi.RegisterHandlersWithOptions(r, server, openapi.GinServerOptions{
		BaseURL:      "/api/v1",
		Middlewares:  nil,
		ErrorHandler: nil,
	})

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", port),
	}
	return s
}

func (s *httpServer) LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return
		}
		token := parts[1]
		user, err := s.UserHandlers.Service.ValidateUser(c, token)
		if err != nil {
			return
		}
		c.Set("user", user)
	}
}
