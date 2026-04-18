// Package handler berisi HTTP handler Gin untuk semua endpoint MNC Wallet.
// Handler parse request, panggil service, dan format response.
package handler

import "github.com/gin-gonic/gin"

// RegisterHealth mendaftarkan endpoint /health untuk smoke test.
// Endpoint ini tidak termasuk 7 endpoint utama.
func RegisterHealth(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "SUCCESS",
			"result": gin.H{"service": "mnc-wallet-api"},
		})
	})
}