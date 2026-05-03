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
	r.GET("/api/places", handlers.GetPlaces)

	r.Run(":8080")
}
