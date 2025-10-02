package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// DB holds the database connection
var DB *sql.DB

// InitDB initializes the database connection and creates tables if they don't exist
func InitDB(dbPath string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	// Create tables if they don't exist
	err = createTables(DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return DB, nil
}

// createTables creates the required database tables
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS cities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		country TEXT NOT NULL,
		subcountry TEXT,
		geonameid TEXT UNIQUE NOT NULL
	);
	CREATE TABLE IF NOT EXISTS city_searches (
		geonameid TEXT PRIMARY KEY,
		search_count INTEGER DEFAULT 1,
		last_searched TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (geonameid) REFERENCES cities(geonameid)
	);
	CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);
	CREATE INDEX IF NOT EXISTS idx_cities_name_lower ON cities(LOWER(name));
	CREATE INDEX IF NOT EXISTS idx_city_searches_count ON city_searches(search_count DESC);
	`
	_, err := db.Exec(query)
	return err
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
