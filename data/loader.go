package data

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"city-autocomplete-api/db"
	"city-autocomplete-api/models"
)

// Cache entry to store results with expiry
type cacheEntry struct {
	results []models.City
	expiry  time.Time
}

// Global cache with thread safety
var (
	searchCache = make(map[string]cacheEntry)
	cacheMutex  = sync.RWMutex{}
)

// LoadCities loads city data from the database
func LoadCities(database *sql.DB) ([]models.City, error) {
	// For this implementation, we're returning all cities from database
	// Though in practice, we might want to limit this or implement pagination
	query := "SELECT name, country, subcountry, geonameid FROM cities ORDER BY name"

	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		err = rows.Scan(&city.Name, &city.Country, &city.Subcountry, &city.GeonameID)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// SearchCities searches for cities based on query string with caching
func SearchCities(database *sql.DB, query string, limit int) ([]models.City, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("%s_%d", query, limit)

	// Check cache first (read lock)
	cacheMutex.RLock()
	if entry, found := searchCache[cacheKey]; found && time.Now().Before(entry.expiry) {
		cacheMutex.RUnlock()
		return entry.results, nil
	}
	cacheMutex.RUnlock()

	// Not in cache, query database
	results, err := db.SearchCities(database, query, limit)
	if err != nil {
		return results, err
	}

	// Add to cache (write lock) - only for successful queries
	cacheMutex.Lock()
	searchCache[cacheKey] = cacheEntry{
		results: results,
		expiry:  time.Now().Add(5 * time.Minute), // Cache for 5 minutes
	}
	cacheMutex.Unlock()

	return results, nil
}

// IncrementSearchCount increments the search count for a city
func IncrementSearchCount(database *sql.DB, geonameid string) error {
	return db.IncrementSearchCount(database, geonameid)
}

// CleanupCache removes expired entries from the cache
func CleanupCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	now := time.Now()
	for key, entry := range searchCache {
		if now.After(entry.expiry) {
			delete(searchCache, key)
		}
	}
}

// StartCacheCleanup starts periodic cleanup of expired cache entries
func StartCacheCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // Run cleanup every minute
		defer ticker.Stop()

		for range ticker.C {
			CleanupCache()
		}
	}()
}
