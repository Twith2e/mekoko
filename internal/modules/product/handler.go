package product

import (
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetProducts(c *gin.Context) {
	var query GetProductsWithFilterQuery
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

	products, count, err := h.Service.GetProductsWithFilter(c.Request.Context(), query.Limit, (query.Page-1)*query.Limit, filter)
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
			BasePrice:          p.BasePrice / 100,
			DiscountPercentage: p.DiscountPercentage,
			Slug:               p.Slug,
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
	product, err := h.Service.GetProductByPublicID(c.Request.Context(), publicID)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	variantDTO := make([]VariantResponse, 0, len(product.Variants))
	for _, v := range product.Variants {
		variantDTO = append(variantDTO, VariantResponse{
			ID:            v.PublicID,
			Color:         v.Color,
			Size:          v.Size,
			ImageURL:      v.ImageURL,
			StockQuantity: v.StockQuantity,
		})
	}

	dto := GetProductsResponse{
		ID:                 product.PublicID,
		Name:               product.Name,
		Description:        product.Description,
		BasePrice:          product.BasePrice,
		DiscountPercentage: product.DiscountPercentage,
		Variants:           variantDTO,
	}

	c.JSON(http.StatusOK, response.APIResponse[GetProductsResponse]{
		Status:  "success",
		Message: "Product fetched successfully",
		Data:    &dto,
	})
}
