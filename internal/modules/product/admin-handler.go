package product

import (
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	Service *Service
}

func NewAdminHandler(service *Service) *AdminHandler {
	return &AdminHandler{Service: service}
}

func (a *AdminHandler) AddProducts(c *gin.Context) {
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

	if err := a.Service.AddProducts(c.Request.Context(), payload); err != nil {
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

func (h *AdminHandler) FetchAllProducts(c *gin.Context) {
	var query GetProducts
	if err := c.ShouldBindQuery(&query); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestQuery)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	products, count, err := h.Service.GetProducts(c.Request.Context(), query.Limit, (query.Page-1)*query.Limit)
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
		product := GetProductsResponse{
			ID:                 p.PublicID,
			Name:               p.Name,
			Description:        p.Description,
			BasePrice:          p.BasePrice,
			DiscountPercentage: p.DiscountPercentage,
			Variants:           []VariantResponse{},
		}
		for _, v := range p.Variants {
			product.Variants = append(product.Variants, VariantResponse{
				ID:            v.PublicID,
				Color:         v.Color,
				Size:          v.Size,
				ImageURL:      v.ImageURL,
				StockQuantity: v.StockQuantity,
			})
		}
		dto = append(dto, product)
	}

	c.JSON(http.StatusOK, response.APIResponse[[]GetProductsResponse]{
		Status: "success",
		Data:   &dto,
		Page:   query.Page,
		Limit:  query.Limit,
		Total:  count,
	})
}
