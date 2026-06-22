package product

import (
	"encoding/json"
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	Service      *Service
	FileUploader FileUploader
}

func NewAdminHandler(service *Service, fileUploader FileUploader) *AdminHandler {
	return &AdminHandler{Service: service, FileUploader: fileUploader}
}

func (a *AdminHandler) AddProducts(c *gin.Context) {
	dataJSON := c.PostForm("data")
	if dataJSON == "" {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}
	imageFiles := form.File["images"]

	var payload AddProductsRequest
	if err := json.Unmarshal([]byte(dataJSON), &payload); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	for i, variant := range payload.Variants {
		if i < len(imageFiles) {
			file, _ := imageFiles[i].Open()
			url, err := a.FileUploader.UploadFile(c.Request.Context(), file, imageFiles[i])
			if err != nil {
				mapped := response.MapError(err)
				c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
					Status: "error",
					Error:  &mapped.Error,
				})
				return
			}
			variant.ImageURL = url
		}
	}

	if err := a.Service.AddProducts(c.Request.Context(), payload); err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: "Product successfully added",
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
