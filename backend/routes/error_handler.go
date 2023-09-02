package routes

import (
	"github.com/gin-gonic/gin"
	"log"
)

func handleError(ctx *gin.Context, err error, code int, errorMsg string) {
	log.Println("Error:", err)
	ctx.JSON(code, gin.H{
		"link":   "error",
		"status": errorMsg,
	})
}
