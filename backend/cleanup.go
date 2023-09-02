package main

import (
	"context"
	"filesharing/config"
	"filesharing/models"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"time"
)

func cleanup(db *gorm.DB, tick time.Duration) {
	for range time.Tick(tick) {
		records, err := dbGetAllFiles(db)
		if err != nil {
			log.Fatalln(err)
		}

		// Check for expired files
		currentTime := time.Now()
		for _, v := range records {
			if currentTime.After(v.Expiration) {
				err := deleteFile(db, config.Env.StorageDir, v)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func dbGetAllFiles(db *gorm.DB) ([]models.File, error) {
	var files []models.File

	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.Env.DbTimeout)
	defer cancel()

	result := db.WithContext(ctxTimeout).Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}

	return files, nil
}

func deleteFile(db *gorm.DB, dir string, file models.File) error {
	// Delete from filesystem
	path := filepath.Join(dir, file.ID)
	err := os.Remove(path)
	if err != nil {
		return err
	}

	// Delete from database
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	result := db.WithContext(ctxTimeout).Delete(file)

	return result.Error
}
