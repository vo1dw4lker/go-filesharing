package routes

import (
	"context"
	"errors"
	"filesharing/config"
	"filesharing/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Upload handles file uploads and creates records in the database.
func Upload(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		expiration, err := extractExpiration(ctx)
		if err != nil {
			handleError(ctx, err, http.StatusBadRequest, "Invalid expiration date")
			return
		}

		file, err := extractFile(ctx)
		if err != nil {
			handleError(ctx, err, http.StatusBadRequest, "File extraction failed")
			return
		}

		fileRecord, err := createFileRecord(db, file, expiration)
		if err != nil {
			handleError(ctx, err, http.StatusInternalServerError, "Failed to create file record")
			return
		}

		err = saveUploadedFile(ctx, file, fileRecord.ID)
		if err != nil {
			handleError(ctx, err, http.StatusInternalServerError, "Failed to save file")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Link":   fileRecord.ID,
			"Status": "OK",
		})
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
	for _, option := range config.Env.AllowedExpOptions {
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

	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.Env.DbTimeout)
	defer cancel()

	if err := db.WithContext(ctxTimeout).Create(fileRecord).Error; err != nil {
		return nil, err
	}

	return fileRecord, nil
}

func saveUploadedFile(ctx *gin.Context, file *multipart.FileHeader, fileName string) error {
	if err := os.MkdirAll(config.Env.StorageDir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(config.Env.StorageDir, fileName)

	err := ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		return err
	}

	return nil
}
