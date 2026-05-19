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

func CreateTrip(c *gin.Context) {
	var trip models.Trip

	// читаем JSON
	if err := c.ShouldBindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// вставка в БД
	query := `
		INSERT INTO trips (title, start_date, end_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := database.DB.QueryRow(
		query,
		trip.Title,
		trip.StartDate,
		trip.EndDate,
	).Scan(&trip.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trip)
}

func UpdateTrip(c *gin.Context) {
	id := c.Param("id")

	var t models.Trip

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	_, err := database.DB.Exec(
		"UPDATE trips SET title=?, start_date=?, end_date=? WHERE id=?",
		t.Title,
		t.StartDate,
		t.EndDate,
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "trip updated",
	})
}
