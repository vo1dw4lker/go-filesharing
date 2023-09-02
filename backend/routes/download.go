package routes

import (
	"filesharing/config"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func Download() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")

		filePath := filepath.Join(config.Env.StorageDir, name)
		ctx.File(filePath)
	}
}
