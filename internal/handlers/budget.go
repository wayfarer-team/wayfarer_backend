package handlers

import (
	"net/http"

	"wayfarer/internal/database"

	"github.com/gin-gonic/gin"
)

func GetBudget(c *gin.Context) {
	id := c.Param("id")

	var total int

	err := database.DB.QueryRow(
		"SELECT COALESCE(SUM(cost), 0) FROM events WHERE trip_id=?",
		id,
	).Scan(&total)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trip_id": id,
		"total":   total,
	})
}
