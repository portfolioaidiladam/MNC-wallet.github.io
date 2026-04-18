package handler

import (
	"net/http"

	"github.com/aidiladam/mnc-wallet/internal/middleware"
	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/service"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/gin-gonic/gin"
)

// TransferHandler menangani endpoint transfer antar-user.
type TransferHandler struct {
	svc service.TransferService
}

// NewTransferHandler membungkus TransferService menjadi Gin handler.
func NewTransferHandler(svc service.TransferService) *TransferHandler {
	return &TransferHandler{svc: svc}
}

// Transfer menangani POST /transfer.
func (h *TransferHandler) Transfer(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Transfer(c.Request.Context(), userID, req.TargetUser, req.Amount, req.Remarks)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusCreated, res)
}