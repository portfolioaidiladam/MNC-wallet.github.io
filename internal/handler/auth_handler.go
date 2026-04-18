package handler

import (
	"errors"
	"net/http"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/service"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/gin-gonic/gin"
)

// AuthHandler menjembatani HTTP request ke AuthService.
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler membungkus AuthService menjadi Gin handler.
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register menangani POST /register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusCreated, res)
}

// Login menangani POST /login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusOK, res)
}

// Refresh menangani POST /refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusOK, res)
}

// mapServiceError memetakan sentinel error service ke HTTP response.
// Internal error (bukan sentinel) dibalas generic 500 supaya tidak bocor.
func mapServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		util.Error(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		util.Error(c, http.StatusUnauthorized, "invalid phone number or pin")
	case errors.Is(err, service.ErrUnauthorized):
		util.Error(c, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, service.ErrNotFound):
		util.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrConflict):
		util.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrInsufficientFunds):
		util.Error(c, http.StatusBadRequest, err.Error())
	default:
		util.Error(c, http.StatusInternalServerError, "internal server error")
	}
}