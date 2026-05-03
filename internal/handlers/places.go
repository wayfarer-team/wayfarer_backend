package handlers

import "github.com/gin-gonic/gin"

func GetPlaces(c *gin.Context) {
	region := c.Query("region")

	places := []gin.H{
		{
			"name":        "Ала-Арча",
			"region":      region,
			"description": "Национальный парк",
		},
	}

	c.JSON(200, places)
}
