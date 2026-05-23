package handlers

import (
	"wayfarer/internal/database"
	"wayfarer/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPlaces(c *gin.Context) {
	region := c.Query("region")

	if region == "" {
		c.JSON(400, gin.H{
			"error": "region is required",
		})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, name, description, image_url, region, category, approximate_cost
		FROM places
		WHERE region=?
	`, region)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	places := []models.Place{}

	for rows.Next() {
		var p models.Place

		rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.ImageURL,
			&p.Region,
			&p.Category,
			&p.ApproximateCost,
		)

		places = append(places, p)
	}

	c.JSON(200, places)
}
