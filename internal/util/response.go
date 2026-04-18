package util

import "github.com/gin-gonic/gin"

// Success membalas response format standar { "status": "SUCCESS", "result": ... }.
func Success(c *gin.Context, status int, result any) {
	c.JSON(status, gin.H{"status": "SUCCESS", "result": result})
}

// Error membalas response format standar { "message": ... }.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"message": message})
}