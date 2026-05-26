package waitlist

import (
	"log"
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

func (h *Handler) JoinWaitlist(c *gin.Context) {
	var payload WaitlistRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	if err := h.service.JoinWaitlist(c.Request.Context(), payload.Email); err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: "You have been added to the waitlist!!!",
	})
}

func (h *Handler) FetchWaitlistedEmails(c *gin.Context) {
	entries, err := h.service.FetchWaitlistedEmails(c.Request.Context())
	if err != nil {
		log.Printf("waitlist: failed to fetch waitlisted emails: %s", err)
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[[]WaitlistEntry]{
		Status:  "success",
		Message: "Waitlisted emails retrieved successfully",
		Data:    &entries,
	})
}

func (h *Handler) GetWaitlistCount(c *gin.Context) {
	count, err := h.service.GetWaitlistCount(c.Request.Context())
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[int]{
		Status:  "success",
		Message: "Waitlist count retrieved successfully",
		Data:    &count,
	})
}
