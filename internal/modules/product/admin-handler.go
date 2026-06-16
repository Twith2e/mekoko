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

func (a *AdminHandler) FetchAllProducts(c *gin.Context) {}
