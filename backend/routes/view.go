package routes

import (
	"context"
	"errors"
	"filesharing/config"
	"filesharing/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// 02 = day, 01 = month, 06 = year
const timeFormatPattern = "02.01.06"

func View(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		record, err := dbGetFile(db, ctx.Param("id"))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"Status": "Not found"})
				return
			}
			handleError(ctx, err, http.StatusInternalServerError, "Failed accessing db")
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Status":         "OK",
			"FileName":       record.FileName,
			"FileSize":       byteToMegabyte(record.FileSize),
			"ExpirationDate": record.Expiration.Format(timeFormatPattern),
			"DownloadLink":   "/api/download/" + record.ID,
		})
	}
}

func dbGetFile(db *gorm.DB, id string) (*models.File, error) {
	record := &models.File{ID: id}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.Env.DbTimeout)
	defer cancel()

	result := db.WithContext(ctxTimeout).First(record)
	if result.Error != nil {
		return nil, result.Error
	}

	return record, nil
}

func byteToMegabyte(bytes int64) string {
	megabytes := float32(bytes) / 1048576.0
	result := fmt.Sprintf("%.2f Mb", megabytes)

	return result
}
