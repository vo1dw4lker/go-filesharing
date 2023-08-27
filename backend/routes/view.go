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
				handleError(ctx, err, http.StatusNotFound)
				return
			}
			handleError(ctx, err, http.StatusInternalServerError)
		}

		ctx.HTML(http.StatusOK, "file-view.html", gin.H{
			"FileName":       record.FileName,
			"FileSize":       record.FileSize,
			"ExpirationDate": record.Expiration.Format(timeFormatPattern),
			"DownloadLink":   "dl", // TODO: return a link to a file here
		})
	}
}

func dbGetFile(db *gorm.DB, id string) (*models.File, error) {
	record := &models.File{ID: id}

	result := db.First(record)
	if result.Error != nil {
		return nil, result.Error
	}

	return record, nil
}

// Set the appropriate headers for downloading
// c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filenameFromDatabase)

// Send the file data to the client
// c.File(filepathToFileOnServer)
