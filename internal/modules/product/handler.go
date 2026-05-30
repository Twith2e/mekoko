package product

import (
	"mekoko/internal/domain"
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddProducts(c *gin.Context) {
	var payload []AddProductsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	if len(payload) <= 0 {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	if err := h.service.AddProducts(c.Request.Context(), payload); err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	message := "Product successfully added"

	if len(payload) > 1 {
		message = "Products successfully added"
	}

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: message,
	})
}

func (h *Handler) GetProducts(c *gin.Context) {
	var query GetProductsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestQuery)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	if query.MaxPrice != nil && query.MinPrice != nil && *query.MaxPrice < *query.MinPrice {
		mapped := response.MapError(appErr.ErrInvalidPriceRange)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	filter := Filter{
		Order:    FilterType(query.Order),
		Color:    query.Color,
		MinPrice: query.MinPrice,
		MaxPrice: query.MaxPrice,
	}

	products, count, err := h.service.GetProducts(c.Request.Context(), query.Limit, (query.Page-1)*query.Limit, filter)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	dto := make([]GetProductsResponse, 0, len(products))
	for _, p := range products {
		prod := GetProductsResponse{
			ID:                 p.PublicID,
			Name:               p.Name,
			Description:        p.Description,
			BasePrice:          p.BasePrice,
			DiscountPercentage: p.DiscountPercentage,
			Variants:           []VariantResponse{},
		}
		for _, v := range p.Variants {
			prod.Variants = append(prod.Variants, VariantResponse{
				ID:            v.PublicID,
				Color:         v.Color,
				Size:          v.Size,
				ImageURL:      v.ImageURL,
				StockQuantity: v.StockQuantity,
			})
		}
		dto = append(dto, prod)
	}

	if query.Page == 0 {
		query.Page = 1
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	c.JSON(http.StatusOK, response.APIResponse[[]GetProductsResponse]{
		Status:  "success",
		Message: "Products fetched successfully",
		Data:    &dto,
		Page:    query.Page,
		Limit:   query.Limit,
		Total:   count,
	})
}

func (h *Handler) GetProductByPublicID(c *gin.Context) {
	publicID := c.Param("public_id")
	product, err := h.service.GetProductByPublicID(c.Request.Context(), publicID)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}
	c.JSON(http.StatusOK, response.APIResponse[*domain.Product]{
		Status:  "success",
		Message: "Product fetched successfully",
		Data:    &product,
	})
}
