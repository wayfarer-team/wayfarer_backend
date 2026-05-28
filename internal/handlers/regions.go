package handlers

import (
	"net/http"
	"wayfarer/internal/database"

	"github.com/gin-gonic/gin"
)

type Region struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

func GetRegions(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, description, image_url FROM regions ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var regions []Region
	for rows.Next() {
		var r Region
		err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		regions = append(regions, r)
	}

	c.JSON(http.StatusOK, regions)
}
