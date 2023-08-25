package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"parcels/models"
)

func main() {
	//initiating tracer
	cleanup := initTracer()
	defer cleanup(context.Background())

	//Setting up SQLite database and connection
	db, err := gorm.Open(sqlite.Open("thedb.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Creating Parcels Table in database
	if err := db.AutoMigrate(&models.Parcel{}); err != nil {
		panic("Failed to connect to database")
	}

	//Adding otelgorm plugging to GORM ORM for db instrumentation
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	//Setting up router
	router := gin.Default()
	//Adding otelgin as middleware to auto-instrument ALL API requests
	router.Use(otelgin.Middleware(serviceName))
	//adding route handlers
	setupRoutes(router, db)

	//Starting the web Server
	if err := router.Run(":8080"); err != nil {
		panic("Failed to start server: ")
	}
}
