package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	"pvz/generated/openapi"
	"pvz/internal/config"
	"time"
)

var _ openapi.ServerInterface = (*httpServer)(nil)

type httpServer struct {
	app application
}

func SetupHTTPServer(app application) *http.Server {
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
	server := httpServer{
		app: app,
	}
	openapi.RegisterHandlers(r, server)
	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", port),
	}
	return s
}

func (h httpServer) PostDummyLogin(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) PostLogin(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) PostProducts(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) GetPvz(c *gin.Context, params openapi.GetPvzParams) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) PostPvz(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) PostPvzPvzIdCloseLastReception(c *gin.Context, pvzId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) PostPvzPvzIdDeleteLastProduct(c *gin.Context, pvzId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) PostReceptions(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h httpServer) PostRegister(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
