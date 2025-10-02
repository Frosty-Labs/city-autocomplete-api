package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"city-autocomplete-api/data"
	"city-autocomplete-api/models"
)

// AutocompleteHandler handles the autocomplete requests
type AutocompleteHandler struct {
	db *sql.DB
}

// NewAutocompleteHandler creates a new instance of AutocompleteHandler
func NewAutocompleteHandler(database *sql.DB) *AutocompleteHandler {
	return &AutocompleteHandler{
		db: database,
	}
}

// Autocomplete handles the HTTP request for city autocomplete
func (h *AutocompleteHandler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	limit := 10 // default limit
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 { // Set a reasonable maximum
				limit = 100
			}
		}
	}

	// Search for cities in the database
	results, err := data.SearchCities(h.db, query, limit)
	if err != nil {
		http.Error(w, "Database error occurred", http.StatusInternalServerError)
		return
	}

	// Use models.City to ensure the import is recognized
	var _ []models.City = results

	// Increment search counts for each returned city in a goroutine to not block the response
	go func() {
		for _, city := range results {
			// We'll ignore the error here as we don't want to affect the response
			_ = data.IncrementSearchCount(h.db, city.GeonameID)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow CORS for web usage
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
