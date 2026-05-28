package handlers

import (
	"database/sql"
	"net/http"
	"wayfarer/internal/database"

	"github.com/gin-gonic/gin"
)

type Attraction struct {
	ID          int     `json:"id"`
	RegionID    int     `json:"region_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Rating      float64 `json:"rating"`
}

func GetAttractions(c *gin.Context) {
	regionID := c.Query("region_id")

	var rows *sql.Rows
	var err error

	if regionID != "" {
		rows, err = database.DB.Query(`
			SELECT id, region_id, title, description, image_url, rating 
			FROM attractions 
			WHERE region_id = ?
			ORDER BY rating DESC
		`, regionID)
	} else {
		rows, err = database.DB.Query(`
			SELECT id, region_id, title, description, image_url, rating 
			FROM attractions 
			ORDER BY rating DESC
		`)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var attractions []Attraction
	for rows.Next() {
		var a Attraction
		err := rows.Scan(&a.ID, &a.RegionID, &a.Title, &a.Description, &a.ImageURL, &a.Rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		attractions = append(attractions, a)
	}

	c.JSON(http.StatusOK, attractions)
}
