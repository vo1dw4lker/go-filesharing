package routes

import (
	"context"
	"errors"
	"filesharing/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Allowed options for file expiration time (in days)
var allowedExpOptions = []int{1, 7, 10}

// Timeout for database operations
const timeout = time.Second * 5
const storageDir = "../storage"

// Upload handles file uploads and creates records in the database.
func Upload(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		expiration, err := extractExpiration(ctx)
		if err != nil {
			handleError(ctx, err, http.StatusBadRequest)
			return
		}

		file, err := extractFile(ctx)
		if err != nil {
			handleError(ctx, err, http.StatusBadRequest)
			return
		}

		fileRecord, err := createFileRecord(db, file, expiration)
		if err != nil {
			handleError(ctx, err, http.StatusInternalServerError)
			return
		}

		err = saveUploadedFile(ctx, file, fileRecord.ID)
		if err != nil {
			handleError(ctx, err, http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"link": fileRecord.ID})
	}
}

func extractExpiration(ctx *gin.Context) (int, error) {
	expStr := ctx.PostForm("exp")
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return 0, err
	}

	if !isExpirationAllowed(exp) {
		return 0, errors.New("expiration not allowed")
	}

	return exp, nil
}

func isExpirationAllowed(exp int) bool {
	for _, option := range allowedExpOptions {
		if exp == option {
			return true
		}
	}
	return false
}

func extractFile(ctx *gin.Context) (*multipart.FileHeader, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func createFileRecord(db *gorm.DB, file *multipart.FileHeader, exp int) (*models.File, error) {
	expDate := time.Now().Add(time.Duration(exp) * time.Hour * 24)

	fileRecord := &models.File{
		ID:         uuid.New().String(),
		FileName:   file.Filename,
		Expiration: expDate,
		FileSize:   file.Size,
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.WithContext(ctxTimeout).Create(fileRecord).Error; err != nil {
		return nil, err
	}

	return fileRecord, nil
}

func saveUploadedFile(ctx *gin.Context, file *multipart.FileHeader, fileName string) error {
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(storageDir, fileName)

	err := ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		return err
	}

	return nil
}

func handleError(ctx *gin.Context, err error, code int) {
	log.Println("Error:", err)
	ctx.Status(code)
}
