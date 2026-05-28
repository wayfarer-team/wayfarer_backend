package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	"wayfarer/internal/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Инициализируем БД
	database.InitDB()
	defer database.CloseDB()

	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Starting server on port %s", port)
	// Используем "0.0.0.0", чтобы слушать все сетевые интерфейсы
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Wayfarer API is running 🚀",
		})
	})

	// ==================== REGIONS ====================
	r.GET("/regions", func(c *gin.Context) {
		rows, err := database.DB.Query("SELECT id, name, description, image_url FROM regions ORDER BY id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var regions []gin.H
		for rows.Next() {
			var id int
			var name, desc, image string
			rows.Scan(&id, &name, &desc, &image)
			regions = append(regions, gin.H{
				"id": id, "name": name, "description": desc, "image_url": image,
			})
		}
		c.JSON(http.StatusOK, regions)
	})

	// ==================== ATTRACTIONS / PLACES ====================
	r.GET("/attractions", func(c *gin.Context) {
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

		var attractions []gin.H
		for rows.Next() {
			var id, regionID int
			var title, desc, image string
			var rating float64
			rows.Scan(&id, &regionID, &title, &desc, &image, &rating)
			attractions = append(attractions, gin.H{
				"id": id, "region_id": regionID, "title": title,
				"description": desc, "image_url": image, "rating": rating,
			})
		}
		c.JSON(http.StatusOK, attractions)
	})

	r.GET("/places", func(c *gin.Context) {
		// Алиас для /attractions
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

		var places []gin.H
		for rows.Next() {
			var id, regionID int
			var title, desc, image string
			var rating float64
			rows.Scan(&id, &regionID, &title, &desc, &image, &rating)
			places = append(places, gin.H{
				"id": id, "region_id": regionID, "title": title,
				"description": desc, "image_url": image, "rating": rating,
			})
		}
		c.JSON(http.StatusOK, places)
	})

	// ==================== TRIPS ====================
	r.GET("/api/trips", func(c *gin.Context) {
		userID := c.Query("user_id")
		var rows *sql.Rows
		var err error

		if userID != "" {
			rows, err = database.DB.Query(`
				SELECT id, user_id, title, description, start_date, end_date, created_at 
				FROM trips 
				WHERE user_id = ?
				ORDER BY created_at DESC
			`, userID)
		} else {
			rows, err = database.DB.Query(`
				SELECT id, user_id, title, description, start_date, end_date, created_at 
				FROM trips 
				ORDER BY created_at DESC
			`)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var trips []gin.H
		for rows.Next() {
			var id, userID int
			var title, description, startDate, endDate, createdAt string
			rows.Scan(&id, &userID, &title, &description, &startDate, &endDate, &createdAt)
			trips = append(trips, gin.H{
				"id": id, "user_id": userID, "title": title, "description": description,
				"start_date": startDate, "end_date": endDate, "created_at": createdAt,
			})
		}
		c.JSON(http.StatusOK, trips)
	})

	r.POST("/api/trips", func(c *gin.Context) {
		var trip struct {
			UserID      int    `json:"user_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			StartDate   string `json:"start_date"`
			EndDate     string `json:"end_date"`
		}

		if err := c.BindJSON(&trip); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := database.DB.Exec(`
			INSERT INTO trips (user_id, title, description, start_date, end_date) 
			VALUES (?, ?, ?, ?, ?)
		`, trip.UserID, trip.Title, trip.Description, trip.StartDate, trip.EndDate)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "Trip created successfully"})
	})

	r.PATCH("/api/trips/:id", func(c *gin.Context) {
		tripID := c.Param("id")

		var updates map[string]interface{}
		if err := c.BindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Строим UPDATE запрос динамически
		query := "UPDATE trips SET "
		args := []interface{}{}
		i := 0
		for key, value := range updates {
			if i > 0 {
				query += ", "
			}
			query += key + " = ?"
			args = append(args, value)
			i++
		}
		query += " WHERE id = ?"
		args = append(args, tripID)

		_, err := database.DB.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Trip updated successfully"})
	})

	// ==================== EVENTS ====================
	r.GET("/api/events", func(c *gin.Context) {
		tripID := c.Query("trip_id")

		rows, err := database.DB.Query(`
			SELECT id, trip_id, day_number, title, start_time, location, cost, description 
			FROM events 
			WHERE trip_id = ?
			ORDER BY day_number, start_time
		`, tripID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var events []gin.H
		for rows.Next() {
			var id, tripID, dayNumber, cost int
			var title, startTime, location, description string
			rows.Scan(&id, &tripID, &dayNumber, &title, &startTime, &location, &cost, &description)
			events = append(events, gin.H{
				"id": id, "trip_id": tripID, "day_number": dayNumber, "title": title,
				"start_time": startTime, "location": location, "cost": cost, "description": description,
			})
		}
		c.JSON(http.StatusOK, events)
	})

	r.POST("/api/events", func(c *gin.Context) {
		var event struct {
			TripID      int    `json:"trip_id"`
			DayNumber   int    `json:"day_number"`
			Title       string `json:"title"`
			StartTime   string `json:"start_time"`
			Location    string `json:"location"`
			Cost        int    `json:"cost"`
			Description string `json:"description"`
		}

		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := database.DB.Exec(`
			INSERT INTO events (trip_id, day_number, title, start_time, location, cost, description) 
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, event.TripID, event.DayNumber, event.Title, event.StartTime, event.Location, event.Cost, event.Description)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "Event created successfully"})
	})

	r.DELETE("/api/events/:id", func(c *gin.Context) {
		eventID := c.Param("id")

		_, err := database.DB.Exec("DELETE FROM events WHERE id = ?", eventID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
	})

	// ==================== AUTH ====================
	r.POST("/api/register", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Хешируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		result, err := database.DB.Exec(`
			INSERT INTO users (username, email, password, created_at) 
			VALUES (?, ?, ?, ?)
		`, user.Username, user.Email, string(hashedPassword), time.Now().Format("2006-01-02 15:04:05"))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username or email already exists"})
			return
		}

		id, _ := result.LastInsertId()
		c.JSON(http.StatusOK, gin.H{
			"id":       id,
			"username": user.Username,
			"email":    user.Email,
			"message":  "User registered successfully",
		})
	})

	r.POST("/api/login", func(c *gin.Context) {
		var login struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user struct {
			ID       int
			Username string
			Email    string
			Password string
		}

		err := database.DB.QueryRow(`
			SELECT id, username, email, password 
			FROM users 
			WHERE email = ?
		`, login.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Проверяем пароль
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"message":  "Login successful",
		})
	})

	// Выводим все зарегистрированные роуты
	log.Println("📋 Registered routes:")
	for _, route := range r.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	log.Println("🚀 Server starting on :8080")
	r.Run(":8080")
}

// package main

// import (
// 	"github.com/gin-contrib/cors"

// 	"wayfarer/internal/handlers"

// 	"github.com/gin-gonic/gin"

// 	"wayfarer/internal/database"
// )

// func main() {
// 	database.InitDB()
// 	r := gin.Default()

// 	r.Use(cors.Default())

// 	// тестовый роут
// 	r.GET("/", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "Wayfarer API is running 🚀",
// 		})
// 	})

// 	// 👇 API
// r.GET("/api/trips", handlers.GetTrips)
// r.POST("/api/trips", handlers.CreateTrip)
// r.PATCH("/api/trips/:id", handlers.UpdateTrip)
// r.GET("/api/places", handlers.GetPlaces)
// r.GET("/api/events", handlers.GetEvents)
// r.POST("/api/events", handlers.CreateEvent)
// r.DELETE("/api/events/:id", handlers.DeleteEvent)
// r.POST("/api/register", handlers.Register)
// r.POST("/api/login", handlers.Login)
// r.GET("/regions", handlers.GetRegions)
// r.GET("/attractions", handlers.GetAttractions)
// r.GET("/places", handlers.GetAttractions)

// 	r.Run(":8080")
// }
