package handler

import (
	"net/http"

	"github.com/aidiladam/mnc-wallet/internal/middleware"
	"github.com/aidiladam/mnc-wallet/internal/service"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/gin-gonic/gin"
)

// TransactionHandler menyajikan GET /transactions.
type TransactionHandler struct {
	svc service.TransactionService
}

// NewTransactionHandler membungkus TransactionService menjadi Gin handler.
func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// List menangani GET /transactions.
func (h *TransactionHandler) List(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	util.Success(c, http.StatusOK, items)
}