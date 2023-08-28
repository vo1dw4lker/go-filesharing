package routes

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func Download() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")

		filePath := filepath.Join(storageDir, name)
		ctx.File(filePath)
	}
}
