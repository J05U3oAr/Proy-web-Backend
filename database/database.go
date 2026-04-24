package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Init opens the SQLite database and creates tables if they don't exist
func Init(path string) {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTables()
	log.Println("Database initialized successfully")
}

func createTables() {
	seriesTable := `
	CREATE TABLE IF NOT EXISTS series (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT    NOT NULL,
		genre       TEXT    NOT NULL DEFAULT '',
		status      TEXT    NOT NULL DEFAULT 'plan_to_watch',
		episodes    INTEGER NOT NULL DEFAULT 0,
		description TEXT    NOT NULL DEFAULT '',
		image_url   TEXT    NOT NULL DEFAULT '',
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(seriesTable)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}
