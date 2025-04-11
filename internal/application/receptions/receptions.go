package receptions

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	"pvz/generated/openapi"
	domain "pvz/internal/domain/receptions"
	"time"
)

type ReceptionService interface {
	CreatePoint(ctx context.Context, point domain.PickUpPoint) error
	GetPointList(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]*domain.PickUpPoint, error)
	CreateReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error)
	AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (domain.Product, error)
	DeleteLastProductFromReception(ctx context.Context, receptionID uuid.UUID) error
}

type ReceptionHandlers struct {
	s ReceptionService
}

func RegisterReceptionHandlers(s ReceptionService) *ReceptionHandlers {
	return &ReceptionHandlers{s}
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
	body := openapi.PostPvzJSONRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID, err := uuid.Parse(body.Id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	point, err := domain.NewPickUpPoint(ID, *body.RegistrationDate, string(body.City))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.s.CreatePoint(c, *point)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, body)
}

func (r *ReceptionHandlers) PostPvzPvzIdCloseLastReception(c *gin.Context, pvzId openapi_types.UUID) {
	reception, err := r.s.CloseLastReception(c, pvzId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convertDomainToResponse(reception))
}

func (r *ReceptionHandlers) PostPvzPvzIdDeleteLastProduct(c *gin.Context, pvzId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (r *ReceptionHandlers) PostReceptions(c *gin.Context) {
	body := openapi.PostReceptionsJSONBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reception, err := r.s.CreateReception(c, body.PvzId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convertDomainToResponse(reception))
}

func convertDomainToResponse(reception domain.Reception) gin.H {
	return gin.H{
		"id":       reception.ID(),
		"dateTime": reception.InitTime(),
		"pvzId":    reception.PvzID(),
		"status":   reception.Status(),
	}
}
