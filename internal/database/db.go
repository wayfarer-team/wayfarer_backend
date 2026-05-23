package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite3", "./wayfarer.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS trips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		start_date TEXT,
		end_date TEXT
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	queryEvents := `
	CREATE TABLE IF NOT EXISTS events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	trip_id INTEGER,
	day_id INTEGER,
	title TEXT,
	start_time TEXT,
	location_name TEXT,
	cost INTEGER,
	description TEXT
);
`

	queryPlaces := `
	CREATE TABLE IF NOT EXISTS places (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		description TEXT,
		image_url TEXT,
		region TEXT,
		category TEXT,
		approximate_cost INTEGER
	);
	`

	_, err = DB.Exec(queryPlaces)

	if err != nil {
		log.Fatal(err)
	}

	DB.Exec(`
	INSERT INTO places(name, description, image_url, region, category, approximate_cost)
	VALUES
	('Ala-Archa', 'National park', 'alaarcha.jpg', 'chuy', 'nature', 200),
	('Burana Tower', 'Historical place', 'burana.jpg', 'chuy', 'museum', 150)
	`)

	queryUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		email TEXT UNIQUE,
		password TEXT
	);
	`

	_, err = DB.Exec(queryUsers)

	if err != nil {
		log.Fatal(err)
	}

	DB.Exec(queryEvents)
}
