package main

import (
	"filesharing/config"
	"filesharing/models"
	"filesharing/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var cfg = config.Env

func main() {
	server := gin.Default()

	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"postgres://%s:%s@db:5432/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
	)))
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.File{})
	if err != nil {
		log.Fatalln(err)
	}

	api := server.Group("/api")
	{
		api.POST("/upload", routes.Upload(db))
		api.GET("/download/:name", routes.Download())
		api.GET("/file/:id", routes.View(db))
	}

	go cleanup(db, time.Hour)

	server.Run(":8081")
}
