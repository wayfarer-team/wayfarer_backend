package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	// Создаем папку для БД если её нет
	dbDir := "internal/database/wayfarer_database"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal("Cannot create database directory:", err)
	}

	// Путь к БД
	dbPath := filepath.Join(dbDir, "wayfarer.db")

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Проверяем подключение
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("✅ Database connected successfully at:", dbPath)

	// Создаем таблицы
	createTables()

	// Добавляем начальные данные
	seedData()
}

func createTables() {
	// Таблица регионов
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS regions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			image_url TEXT
		)
	`)
	if err != nil {
		log.Fatal("Failed to create regions table:", err)
	}

	// Таблица достопримечательностей
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS attractions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			region_id INTEGER,
			title TEXT NOT NULL,
			description TEXT,
			image_url TEXT,
			rating REAL DEFAULT 0,
			FOREIGN KEY (region_id) REFERENCES regions(id)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create attractions table:", err)
	}

	// Таблица пользователей
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	// Таблица поездок
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS trips (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT NOT NULL,
			description TEXT,
			start_date TEXT,
			end_date TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create trips table:", err)
	}

	// Таблица событий в поездках
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			trip_id INTEGER,
			day_number INTEGER,
			title TEXT NOT NULL,
			start_time TEXT,
			location TEXT,
			cost INTEGER DEFAULT 0,
			description TEXT,
			FOREIGN KEY (trip_id) REFERENCES trips(id)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create events table:", err)
	}

	log.Println("✅ All tables created successfully")
}

func seedData() {
	// Проверяем, есть ли уже данные
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM regions").Scan(&count)
	if err != nil {
		log.Println("Warning: Could not check regions count:", err)
		return
	}

	if count > 0 {
		log.Println("Database already has data, skipping seed")
		return
	}

	log.Println("Seeding initial data...")

	// Добавляем регионы
	regions := []struct {
		name, desc, image string
	}{
		{"Бишкек", "Столица Кыргызстана и культурный центр страны", "/images/bishkek.jpg"},
		{"Ош", "Древний город на юге Кыргызстана с 3000-летней историей", "/images/osh.jpg"},
		{"Иссык-Куль", "Известное высокогорное незамерзающее озеро Кыргызстана", "/images/issyk_kul.jpg"},
		{"Нарын", "Высокогорный регион первозданных озер и суровой природы", "/images/naryn.jpg"},
		{"Талас", "Родина легендарного героя Манаса и сакральных памятников", "/images/talas.jpg"},
		{"Чуй", "Плодородная долина, соединяющая столицу с великими ущельями", "/images/chui.jpg"},
		{"Джалал-Абад", "Край реликтовых ореховых лесов и биосферных заповедников", "/images/jalal_abad.jpg"},
		{"Баткен", "Уникальный регион суровых гранитных пиков и древних пещер", "/images/batken.jpg"},
	}

	for _, r := range regions {
		_, err := DB.Exec(`
			INSERT INTO regions (name, description, image_url) 
			VALUES (?, ?, ?)
		`, r.name, r.desc, r.image)
		if err != nil {
			log.Println("Error seeding region:", err)
		}
	}

	// Добавляем достопримечательности
	attractions := []struct {
		regionID    int
		title       string
		description string
		imageURL    string
		rating      float64
	}{
		{6, "Ущелье Ала-Арча", "Популярное место с ледниками и бурными реками вблизи столицы", "/images/ala_archa.jpg", 4.8},
		{6, "Башня Бурана", "Древний минарет XI века и археологический музей под открытым небом", "/images/burana.jpg", 4.7},
		{3, "Горячие источники Ак-Суу", "Термальные источники среди горной природы", "/images/ak_suu.jpg", 4.6},
		{3, "Ущелье Барскоон", "Известное ущелье с каскадными водопадами", "/images/barskoon.jpg", 4.5},
		{1, "Площадь Ала-Тоо", "Главная площадь столицы с памятниками и фонтанами", "/images/ala_too.jpg", 4.9},
		{2, "Гора Сулайман-Тоо", "Священная гора и объект Всемирного наследия ЮНЕСКО", "/images/sulayman_too.jpg", 4.8},
	}

	for _, a := range attractions {
		_, err := DB.Exec(`
			INSERT INTO attractions (region_id, title, description, image_url, rating) 
			VALUES (?, ?, ?, ?, ?)
		`, a.regionID, a.title, a.description, a.imageURL, a.rating)
		if err != nil {
			log.Println("Error seeding attraction:", err)
		}
	}

	log.Println("✅ Seed data added successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
