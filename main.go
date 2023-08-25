package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"parcels/models"
)

func main() {

	//Setting up SQLite database and connection
	db, err := gorm.Open(sqlite.Open("thedb.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Creating PArcels Table in database
	if err := db.AutoMigrate(&models.Parcel{}); err != nil {
		panic("Failed to connect to database")
	}

	//Setting up router
	router := gin.Default()
	setupRoutes(router, db)

	//Starting Server
	if err := router.Run(":8080"); err != nil {
		panic("Failed to start server: ")
	}
}
