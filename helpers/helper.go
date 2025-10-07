package helpers

import (
	"github.com/gin-gonic/gin"
)

// Using for sending error or success message
func SendMessageAsJson(ctx *gin.Context, message string, statusCode int) {
	ctx.JSON(statusCode, gin.H{"message": message})
}
