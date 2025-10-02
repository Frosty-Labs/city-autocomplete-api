package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"city-autocomplete-api/data"
	"city-autocomplete-api/db"
	"city-autocomplete-api/handlers"
)

const defaultPort = "8080"
const defaultDBPath = "cities.db"
const defaultCSVPath = "world-cities.csv"

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// Initialize database
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.CloseDB()

	// Check if cities table is empty and populate from CSV if needed
	populated, err := checkAndPopulateDB(database)
	if err != nil {
		log.Fatalf("Error checking/populating database: %v", err)
	}

	if populated {
		log.Println("Database populated from CSV file")
	} else {
		log.Println("Using existing database")
	}

	// Count cities in database
	var cityCount int
	err = database.QueryRow("SELECT COUNT(*) FROM cities").Scan(&cityCount)
	if err != nil {
		log.Printf("Warning: Could not count cities: %v", err)
	} else {
		log.Printf("Database contains %d cities", cityCount)
	}

	// Start cache cleanup goroutine
	data.StartCacheCleanup()

	// Create an instance of the autocomplete handler with database connection
	autoHandler := handlers.NewAutocompleteHandler(database)

	// Set up HTTP routes
	http.HandleFunc("/autocomplete", autoHandler.Autocomplete)

	// Serve static files (optional, for a simple frontend)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "City Autocomplete API\n\nUse /autocomplete?q=searchterm to search for cities")
	})

	// Determine port to run on
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	fmt.Printf("Starting server on port %s\n", port)
	fmt.Printf("API endpoint: http://localhost:%s/autocomplete?q=ber&limit=5\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// checkAndPopulateDB checks if the database has cities and populates it from CSV if empty
func checkAndPopulateDB(database *sql.DB) (bool, error) {
	exists, err := db.CheckIfCitiesExist(database)
	if err != nil {
		return false, err
	}

	if !exists {
		// Get CSV path from environment or use default
		csvPath := os.Getenv("CSV_PATH")
		if csvPath == "" {
			csvPath = defaultCSVPath
		}

		log.Printf("Populating database from CSV file: %s", csvPath)
		err = db.PopulateCitiesFromCSV(database, csvPath)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
