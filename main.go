package main

import (
	"wayfarer/internal/handlers"

	"github.com/gin-gonic/gin"

	"wayfarer/internal/database"
)

func main() {
	database.InitDB()
	r := gin.Default()

	// тестовый роут
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Wayfarer API is running 🚀",
		})
	})

	// 👇 API
	r.GET("/api/trips", handlers.GetTrips)
	r.POST("/api/trips", handlers.CreateTrip)
	r.PATCH("/api/trips/:id", handlers.UpdateTrip)
	r.GET("/api/places", handlers.GetPlaces)
	r.GET("/api/events", handlers.GetEvents)
	r.POST("/api/events", handlers.CreateEvent)
	r.DELETE("/api/events/:id", handlers.DeleteEvent)

	r.Run(":8080")
}
