package receptions

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	"pvz/generated/openapi"
	domain "pvz/internal/domain/receptions"
	"pvz/internal/domain/users"
	"time"
)

type ReceptionService interface {
	CreatePoint(ctx context.Context, point domain.PickUpPoint) error
	GetPointList(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]*domain.PickUpPoint, error)
	CreateReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (domain.Reception, error)
	AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (domain.Product, uuid.UUID, error)
	DeleteLastProductFromReception(ctx context.Context, pvzID uuid.UUID) error
}

type ReceptionHandlers struct {
	Service ReceptionService
}

func RegisterReceptionHandlers(s ReceptionService) *ReceptionHandlers {
	return &ReceptionHandlers{s}
}

func (r *ReceptionHandlers) PostProducts(c *gin.Context) {
	role, err := getRole(c)
	if err != nil || role != users.Employee {
		c.JSON(http.StatusForbidden, gin.H{"error": "only employees can add products"})
		return
	}

	body := openapi.PostProductsJSONBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID, err := uuid.Parse(body.PvzId.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, receptionID, err := r.Service.AddProduct(c, ID, string(body.Type))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convertProductToResponse(product, receptionID.String()))
}

func (r *ReceptionHandlers) GetPvz(c *gin.Context, params openapi.GetPvzParams) {
	role, err := getRole(c)
	if err != nil || (role != users.Employee && role != users.Moderator) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only employees can add products"})
		return
	}
	if params.Page == nil || params.Limit == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}
	if params.EndDate == nil {
		tmp := time.Now()
		params.EndDate = &tmp
	}
	if params.StartDate == nil {
		tmp := time.Now().AddDate(-1, -1, 0)
		params.StartDate = &tmp
	}
	if params.EndDate.Before(*params.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time window"})
		return
	}
	points, err := r.Service.GetPointList(c, *params.StartDate, *params.EndDate, *params.Page, *params.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convertPointListToResponse(points))
}

func (r *ReceptionHandlers) PostPvz(c *gin.Context) {
	role, err := getRole(c)
	if err != nil || role != users.Moderator {
		c.JSON(http.StatusForbidden, gin.H{"error": "only moderator can add points"})
		return
	}
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
	point, err := domain.NewPickUpPoint(ID, *body.RegistrationDate, string(body.City), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.Service.CreatePoint(c, *point)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, body)
}

func (r *ReceptionHandlers) PostPvzPvzIdCloseLastReception(c *gin.Context, pvzId openapi_types.UUID) {
	role, err := getRole(c)
	if err != nil || role != users.Employee {
		c.JSON(http.StatusForbidden, gin.H{"error": "only employees can close receptions"})
		return
	}
	reception, err := r.Service.CloseLastReception(c, pvzId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convertReceptionToResponse(reception, pvzId.String()))
}

func (r *ReceptionHandlers) PostPvzPvzIdDeleteLastProduct(c *gin.Context, pvzId openapi_types.UUID) {
	role, err := getRole(c)
	if err != nil || role != users.Employee {
		c.JSON(http.StatusForbidden, gin.H{"error": "only employees can delete products"})
		return
	}
	err = r.Service.DeleteLastProductFromReception(c, pvzId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (r *ReceptionHandlers) PostReceptions(c *gin.Context) {
	role, err := getRole(c)
	if err != nil || role != users.Employee {
		c.JSON(http.StatusForbidden, gin.H{"error": "only employees can open receptions"})
		return
	}
	body := openapi.PostReceptionsJSONBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.PvzId.String() == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}

	reception, err := r.Service.CreateReception(c, body.PvzId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convertReceptionToResponse(reception, body.PvzId.String()))
}

func convertReceptionToResponse(reception domain.Reception, pvzId string) gin.H {
	result := gin.H{
		"id":       reception.ID(),
		"dateTime": reception.InitTime(),
		"pvzId":    pvzId,
		"status":   reception.Status(),
	}
	if len(reception.Products()) > 0 {
		products := make([]gin.H, 0, len(reception.Products()))
		for _, product := range reception.Products() {
			products = append(products, convertProductToResponse(product, reception.ID().String()))
		}
		result["products"] = products
	}
	return result
}

func convertProductToResponse(product domain.Product, receptionID string) gin.H {
	return gin.H{
		"id":          product.ID(),
		"dateTime":    product.ArrivalTime(),
		"type":        product.ProductType(),
		"receptionId": receptionID,
	}
}

func convertPointListToResponse(points []*domain.PickUpPoint) []gin.H {
	var result []gin.H
	for _, pickUpPoint := range points {
		point := gin.H{
			"id":               pickUpPoint.ID(),
			"registrationDate": pickUpPoint.RegistrationDate(),
			"city":             pickUpPoint.City(),
		}
		receptions := make([]gin.H, 0, len(pickUpPoint.Receptions()))
		for _, r := range pickUpPoint.Receptions() {
			receptions = append(receptions, convertReceptionToResponse(r, pickUpPoint.ID().String()))
		}
		result = append(result, gin.H{
			"pvz":        point,
			"receptions": receptions,
		})
	}
	return result
}

func getRole(c *gin.Context) (string, error) {
	var user users.User
	userV, ok := c.Get("user")
	if !ok {
		return "", errors.New("user not found")
	}
	if user, ok = userV.(users.User); !ok {
		return "", errors.New("user not found")
	}
	return user.Role(), nil
}
