package auth

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

func (h *Handler) Register(c *gin.Context) {
	var req RegistrationRequest
	if err := c.ShouldBind(&req); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	uat, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	dto := RegistrationResponse{
		ID:           uat.User.UUID,
		FirstName:    uat.User.FirstName,
		Email:        uat.User.Email,
		AccessToken:  uat.Tokens.AccessToken,
		RefreshToken: uat.Tokens.RefreshToken,
	}

	c.JSON(http.StatusOK, response.APIResponse[RegistrationResponse]{
		Status:  "success",
		Message: "Registration was successful",
		Data:    &dto,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[LoginResponse]{
		Status:  "success",
		Message: "Logged in",
		Data:    (*LoginResponse)(tokens),
	})
}

func (h *Handler) Logout(c *gin.Context) {
	pid := strings.TrimSpace(c.GetString(middleware.PublicIDContextKey))
	if pid == "" {
		mapped := response.MapError(appErr.ErrUnauthorized)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	sid := strings.TrimSpace(c.GetString(middleware.SessionIDContextKey))
	if sid == "" {
		mapped := response.MapError(appErr.ErrInvalidSession)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: "Logged out",
	})
}
