package auth

import (
	appErr "mekoko/internal/errors"
	"mekoko/internal/middleware"
	"mekoko/internal/response"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	isProd  bool
}

func NewHandler(service *Service, isProd bool) *Handler {
	return &Handler{service: service, isProd: isProd}
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

	c.SetCookie(CookieName, uat.Tokens.RefreshToken, int(time.Until(uat.Tokens.ExpiresAt).Seconds()), "/api/v1/auth", "", h.isProd, true)

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

	c.SetCookie(CookieName, tokens.RefreshToken, int(time.Until(tokens.ExpiresAt).Seconds()), "/api/v1/auth", "", h.isProd, true)

	c.JSON(http.StatusOK, response.APIResponse[LoginResponse]{
		Status:  "success",
		Message: "Logged in",
		Data: &LoginResponse{
			AccessToken: tokens.AccessToken,
		},
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	pid := strings.TrimSpace(c.GetString(middleware.PublicIDContextKey))
	if pid == "" {
		mapped := response.MapError(appErr.ErrUnauthorized)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	var payload PasswordChangeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), pid, payload); err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: "Password change was successful",
	})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var payload ForgotPasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), payload); err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: "If that email exists, you'll receive a reset link",
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var payload ResetPasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		mapped := response.MapError(appErr.ErrInvalidRequestBody)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), payload); err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: "Password reset was successful",
	})
}

func (h *Handler) RefreshAccessToken(c *gin.Context) {
	cookie, err := (c.Request.Cookie(CookieName))
	if err != nil {
		mapped := response.MapError(appErr.ErrInvalidSession)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	rToken := strings.TrimSpace(cookie.Value)
	if rToken == "" {
		mapped := response.MapError(appErr.ErrInvalidSession)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	tokens, err := h.service.RefreshAccessToken(c.Request.Context(), rToken)
	if err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.SetCookie(CookieName, tokens.RefreshToken, int(time.Until(tokens.ExpiresAt).Seconds()), "/api/v1/auth", "", h.isProd, true)

	dto := LoginResponse{
		AccessToken: tokens.AccessToken,
	}

	c.JSON(http.StatusOK, response.APIResponse[LoginResponse]{
		Status:  "success",
		Message: "Access token refreshed successfully",
		Data:    &dto,
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

	if err := h.service.Logout(c.Request.Context(), sid); err != nil {
		mapped := response.MapError(err)
		c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
			Status: "error",
			Error:  &mapped.Error,
		})
		return
	}

	c.SetCookie(CookieName, "", -1, "/api/v1/auth", "", h.isProd, true)

	c.JSON(http.StatusOK, response.APIResponse[any]{
		Status:  "success",
		Message: "Logged out",
	})
}
