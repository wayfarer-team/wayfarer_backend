package handlers

import (
	"net/http"

	"wayfarer/internal/database"
	"wayfarer/internal/models"

	"github.com/gin-gonic/gin"
)

// GET /api/events
func GetEvents(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, day_id, title, start_time, location_name, cost, description 
		FROM events
	`)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	events := []models.Event{}

	for rows.Next() {
		var e models.Event

		rows.Scan(
			&e.ID,
			&e.DayID,
			&e.Title,
			&e.StartTime,
			&e.LocationName,
			&e.Cost,
			&e.Description,
		)

		events = append(events, e)
	}

	c.JSON(http.StatusOK, events)
}

// POST /api/events
func CreateEvent(c *gin.Context) {
	var e models.Event

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO events(trip_id, day_id, title, start_time, location_name, cost, description)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		e.TripID,
		e.DayID,
		e.Title,
		e.StartTime,
		e.LocationName,
		e.Cost,
		e.Description,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "event created",
	})
}

// DELETE /api/events/:id
func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	_, err := database.DB.Exec(
		"DELETE FROM events WHERE id=?",
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "event deleted",
	})
}
