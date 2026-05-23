package handlers

import (
	"wayfarer/internal/database"
	"wayfarer/internal/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid data",
		})
		return
	}

	_, err := database.DB.Exec(
		"INSERT INTO users(username, email, password) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.Password,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "user registered",
	})
}

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid data",
		})
		return
	}

	var dbUser models.User

	err := database.DB.QueryRow(
		"SELECT id, username, email, password FROM users WHERE email=?",
		user.Email,
	).Scan(
		&dbUser.ID,
		&dbUser.Username,
		&dbUser.Email,
		&dbUser.Password,
	)

	if err != nil {
		c.JSON(401, gin.H{
			"error": "user not found",
		})
		return
	}

	if dbUser.Password != user.Password {
		c.JSON(401, gin.H{
			"error": "wrong password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "login successful",
		"user":    dbUser,
	})
}
