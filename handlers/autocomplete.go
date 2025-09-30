package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"city-autocomplete-api/models"
)

// AutocompleteHandler handles the autocomplete requests
type AutocompleteHandler struct {
	cities []models.City
}

// NewAutocompleteHandler creates a new instance of AutocompleteHandler
func NewAutocompleteHandler(cities []models.City) *AutocompleteHandler {
	return &AutocompleteHandler{
		cities: cities,
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

	query = strings.ToLower(query)
	var results []models.City

	// First, find cities where the query matches the beginning of the name (higher priority)
	var prefixMatches []models.City
	var substringMatches []models.City

	for _, city := range h.cities {
		cityNameLower := strings.ToLower(city.Name)
		if strings.HasPrefix(cityNameLower, query) {
			prefixMatches = append(prefixMatches, city)
		} else if strings.Contains(cityNameLower, query) {
			substringMatches = append(substringMatches, city)
		}

		// Stop early if we have enough results
		if len(prefixMatches)+len(substringMatches) >= limit {
			break
		}
	}

	// Combine results with prefix matches first (higher priority)
	results = append(prefixMatches, substringMatches...)

	// If we didn't get enough results, continue searching
	if len(results) < limit {
		for _, city := range h.cities {
			cityNameLower := strings.ToLower(city.Name)
			alreadyAdded := false
			// Check if already in results
			for _, res := range results {
				if res.GeonameID == city.GeonameID {
					alreadyAdded = true
					break
				}
			}

			if !alreadyAdded && strings.Contains(cityNameLower, query) {
				results = append(results, city)
				if len(results) >= limit {
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow CORS for web usage
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
