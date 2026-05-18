package order

import (
	appErr "mekoko/internal/errors"
	"mekoko/internal/middleware"
	"mekoko/internal/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	pid := strings.TrimSpace(c.GetString(middleware.PublicIDContextKey))
	if pid == "" {
		mapped := response.MapError(appErr.ErrUnauthorized)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	var payload CreateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), pid, payload)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	dto := CreateOrderResponse{
		OrderID: order.PublicID,
	}

	c.JSON(http.StatusOK, response.APIResponse[CreateOrderResponse]{
		Status:  "success",
		Message: "Order created",
		Data:    &dto,
	})
}
