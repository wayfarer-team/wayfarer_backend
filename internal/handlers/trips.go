package handlers

import (
	"net/http"

	"wayfarer/internal/database"
	"wayfarer/internal/models"

	"github.com/gin-gonic/gin"
)

func GetTrips(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, start_date, end_date FROM trips")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	trips := []models.Trip{}

	for rows.Next() {
		var t models.Trip

		err := rows.Scan(&t.ID, &t.Title, &t.StartDate, &t.EndDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		trips = append(trips, t)
	}

	c.JSON(http.StatusOK, trips)
}
