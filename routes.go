package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"parcels/models"
	"time"
)

func setupRoutes(router *gin.Engine, db *gorm.DB) {
	// Create a new parcel
	router.POST("/parcels", func(c *gin.Context) {
		// Simulate authentication delay
		time.Sleep(time.Millisecond * time.Duration(randomDelay(500, 800)))

		var parcel models.Parcel
		c.BindJSON(&parcel)

		// Simulate external API call delay
		time.Sleep(time.Millisecond * time.Duration(randomDelay(400, 800)))
		//generate fees from external service
		if err := parcel.GenerateFee(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate fee:"})
			return
		}
		//inserting record into db
		if err := db.Create(&parcel).Error; err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to insert record"})
			return
		}

		//returning response
		c.JSON(http.StatusOK, parcel)
	})

	// Retrieve all parcels
	router.GET("/parcels", func(c *gin.Context) {
		// Simulate authentication delay
		time.Sleep(time.Millisecond * time.Duration(randomDelay(500, 800)))

		var parcels []models.Parcel

		//query bd for records
		if err := db.Find(&parcels).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		}

		c.JSON(http.StatusOK, parcels)
	})

	// Retrieve a specific parcel with ID from URL
	router.GET("/parcels/:id", func(c *gin.Context) {
		// Simulate authentication delay
		time.Sleep(time.Millisecond * time.Duration(randomDelay(500, 800)))

		var parcel models.Parcel
		//extracting ID from URL
		parcel.ID = c.Param("id")

		//querying DB
		if err := db.First(&parcel).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		}

		//returning response
		c.JSON(http.StatusOK, parcel)
	})

	// Update a parcel
	router.PUT("/parcels/:id", func(c *gin.Context) {
		// Simulate authentication delay
		time.Sleep(time.Millisecond * time.Duration(randomDelay(500, 800)))

		var updatedParcel models.Parcel

		// reading data into parcel
		if err := c.BindJSON(&updatedParcel); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to understand payload"})
			return
		}

		var parcel models.Parcel
		parcel.ID = c.Param("id")
		//retrieving original record

		if err := db.First(&parcel).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parcel not found"})
			return
		}

		// Simulate external API call delay
		time.Sleep(time.Millisecond * time.Duration(randomDelay(400, 800)))
		//generate fees from external service
		if err := parcel.GenerateFee(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate fee"})
			return
		}

		//Saving record to database
		if err := db.Model(&parcel).Updates(&updatedParcel).Error; err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to update record"})
			return
		}
		//returning response
		c.JSON(http.StatusOK, updatedParcel)
	})

	// Delete a parcel
	router.DELETE("/parcels/:id", func(c *gin.Context) {
		// Simulate authentication delay
		time.Sleep(time.Millisecond * time.Duration(randomDelay(300, 600)))

		var parcel models.Parcel
		parcel.ID = c.Param("id")
		//db.First(&parcel)

		if err := db.Delete(&parcel).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete. Record not found"})
		}

		//returning response
		c.JSON(http.StatusNoContent, nil)
	})
}

// random delay time
func randomDelay(min, max int) int {
	return rand.Intn(max-min) + min
}
