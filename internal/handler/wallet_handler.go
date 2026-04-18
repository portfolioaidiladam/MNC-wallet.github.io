package handler

import (
	"net/http"

	"github.com/aidiladam/mnc-wallet/internal/middleware"
	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/service"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/gin-gonic/gin"
)

// WalletHandler menangani endpoint topup dan payment.
type WalletHandler struct {
	svc service.WalletService
}

// NewWalletHandler membungkus WalletService menjadi Gin handler.
func NewWalletHandler(svc service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

// TopUp menangani POST /topup.
func (h *WalletHandler) TopUp(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.TopUp(c.Request.Context(), userID, req.Amount)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusCreated, res)
}

// Pay menangani POST /pay.
func (h *WalletHandler) Pay(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Pay(c.Request.Context(), userID, req.Amount, req.Remarks)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusCreated, res)
}