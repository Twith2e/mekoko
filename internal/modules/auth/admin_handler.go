package auth

import (
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	Service *Service
	IsProd  bool
}

func NewAdminHandler(service *Service, isProd bool) *AdminHandler {
	return &AdminHandler{
		Service: service,
		IsProd:  isProd,
	}
}

func (h *AdminHandler) Register(c *gin.Context) {
	var req RegistrationRequest
	if err := c.ShouldBind(&req); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	uat, err := h.Service.Register(c.Request.Context(), req, AdminRole)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.SetCookie(CookieName, uat.Tokens.RefreshToken, int(time.Until(uat.Tokens.ExpiresAt).Seconds()), "/api/v1/auth", "", h.IsProd, true)

	dto := RegistrationResponse{
		ID:          uat.User.UUID,
		FirstName:   uat.User.FirstName,
		Email:       uat.User.Email,
		AccessToken: uat.Tokens.AccessToken,
	}

	c.JSON(http.StatusOK, response.APIResponse[RegistrationResponse]{
		Status:  "success",
		Message: "Registration was successful",
		Data:    &dto,
	})
}
