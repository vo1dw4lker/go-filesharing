package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// TODO: read constants below from config
// Timeout for database operations
const dbTimeout = time.Second * 5

// Allowed options for file expiration time (in days)
var allowedExpOptions = []int{1, 7, 10}

const storageDir = "../storage"

func handleError(ctx *gin.Context, err error, code int, errorMsg string) {
	log.Println("Error:", err)
	ctx.JSON(code, gin.H{
		"link":   "error",
		"status": errorMsg,
	})
}
