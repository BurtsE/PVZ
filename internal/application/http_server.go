package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz/generated/openapi"
	"pvz/internal/application/receptions"
	"pvz/internal/application/users"
	"pvz/internal/config"
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
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\" \"body size\": %d\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
			param.Request.ContentLength,
		)
	}))
	server := &httpServer{
		UserHandlers:      users.RegisterUserHandlers(userService),
		ReceptionHandlers: &receptions.ReceptionHandlers{},
	}
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
