package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"city-autocomplete-api/data"
	"city-autocomplete-api/handlers"
)

const defaultPort = "8080"

func main() {
	// Load city data from CSV file
	cities, err := data.LoadCities("world-cities.csv")
	if err != nil {
		log.Fatalf("Error loading city data: %v", err)
	}

	fmt.Printf("Loaded %d cities\n", len(cities))

	// Create an instance of the autocomplete handler
	autoHandler := handlers.NewAutocompleteHandler(cities)

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
