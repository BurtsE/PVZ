package receptions

import (
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"pvz/generated/openapi"
)

type ReceptionService interface {
}

type ReceptionHandlers struct {
	s ReceptionService
}

func (r *ReceptionHandlers) PostProducts(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *ReceptionHandlers) GetPvz(c *gin.Context, params openapi.GetPvzParams) {
	//TODO implement me
	panic("implement me")
}

func (r *ReceptionHandlers) PostPvz(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *ReceptionHandlers) PostPvzPvzIdCloseLastReception(c *gin.Context, pvzId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (r *ReceptionHandlers) PostPvzPvzIdDeleteLastProduct(c *gin.Context, pvzId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (r *ReceptionHandlers) PostReceptions(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
