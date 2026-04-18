package handler

import (
	"net/http"

	"github.com/aidiladam/mnc-wallet/internal/middleware"
	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/service"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/gin-gonic/gin"
)

// ProfileHandler menyajikan GET /profile dan PUT /profile.
type ProfileHandler struct {
	svc service.ProfileService
}

// NewProfileHandler membungkus ProfileService menjadi Gin handler.
func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// Get menangani GET /profile.
func (h *ProfileHandler) Get(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	res, err := h.svc.Get(c.Request.Context(), userID)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusOK, res)
}

// Update menangani PUT /profile.
func (h *ProfileHandler) Update(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Update(c.Request.Context(), userID, &req)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusOK, res)
}