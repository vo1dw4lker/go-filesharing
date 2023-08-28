package routes

import (
	"errors"
	"filesharing/models"
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
				handleError(ctx, err, http.StatusNotFound, "File not found")
				return
			}
			handleError(ctx, err, http.StatusInternalServerError, "Failed accessing db")
		}

		ctx.HTML(http.StatusOK, "file-view.html", gin.H{
			"FileName":       record.FileName,
			"FileSize":       record.FileSize, // TODO: convert to megabytes/kilobytes
			"ExpirationDate": record.Expiration.Format(timeFormatPattern),
			"DownloadLink":   "/api/download/" + record.ID,
		})
	}
}

func dbGetFile(db *gorm.DB, id string) (*models.File, error) {
	record := &models.File{ID: id}

	// TODO: use context
	result := db.First(record)
	if result.Error != nil {
		return nil, result.Error
	}

	return record, nil
}
