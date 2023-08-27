package main

import (
	"filesharing/models"
	"filesharing/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("../frontend/templates/*")

	db, err := gorm.Open(postgres.Open("postgres://postgres:passwd@db:5432/filesharing"))
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
	}

	server.GET("/file/:id", routes.View(db))

	server.Run(":8081")
}
