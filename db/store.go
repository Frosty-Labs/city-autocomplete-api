package db

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"city-autocomplete-api/models"
)

// PopulateCitiesFromCSV populates the database with cities from a CSV file
func PopulateCitiesFromCSV(db *sql.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to comma

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Prepare the insert statement
	insertStmt, err := db.Prepare("INSERT INTO cities (name, country, subcountry, geonameid) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer insertStmt.Close()

	// Skip header row if it exists
	header := true
	for i, record := range records {
		if header {
			header = false
			continue
		}

		if len(record) >= 4 {
			// Clean up the fields by trimming spaces
			name := strings.TrimSpace(record[0])
			country := strings.TrimSpace(record[1])
			subcountry := strings.TrimSpace(record[2])
			geonameid := strings.TrimSpace(record[3])

			_, err = insertStmt.Exec(name, country, subcountry, geonameid)
			if err != nil {
				return fmt.Errorf("failed to insert city at row %d: %v", i, err)
			}
		}
	}

	return nil
}

// CheckIfCitiesExist checks if the cities table has any records
func CheckIfCitiesExist(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM cities").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SearchCities searches for cities based on the query string
func SearchCities(db *sql.DB, query string, limit int) ([]models.City, error) {
	// Get prefix and substring queries
	prefixQuery := query + "%"
	substringQuery := "%" + query + "%"

	// Query for both prefix and substring matches
	// Order by: prefix match priority, then popularity (search count), then city name
	queryStmt := `
		SELECT c.name, c.country, c.subcountry, c.geonameid
		FROM cities c
		LEFT JOIN city_searches cs ON c.geonameid = cs.geonameid
		WHERE LOWER(c.name) LIKE LOWER(?) ESCAPE '\'
		ORDER BY 
			CASE 
				WHEN LOWER(c.name) LIKE LOWER(?) ESCAPE '\' THEN 1  -- prefix match first
				ELSE 2  -- substring match second
			END,
			COALESCE(cs.search_count, 0) DESC,  -- popularity: higher counts first
			c.name
		LIMIT ?`

	stmt, err := db.Prepare(queryStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(substringQuery, prefixQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.City
	for rows.Next() && len(results) < limit {
		var city models.City
		err = rows.Scan(&city.Name, &city.Country, &city.Subcountry, &city.GeonameID)
		if err != nil {
			return nil, err
		}
		results = append(results, city)
	}

	return results, nil
}

// IncrementSearchCount increments the search count for a city
func IncrementSearchCount(db *sql.DB, geonameid string) error {
	query := `
		INSERT INTO city_searches (geonameid, search_count, last_searched) 
		VALUES (?, 1, datetime('now')) 
		ON CONFLICT(geonameid) 
		DO UPDATE SET 
			search_count = search_count + 1,
			last_searched = datetime('now')
		WHERE geonameid = ?`
	_, err := db.Exec(query, geonameid, geonameid)
	return err
}

// GetPopularityScore returns the search count for a city
func GetPopularityScore(db *sql.DB, geonameid string) (int, error) {
	var count int
	query := "SELECT search_count FROM city_searches WHERE geonameid = ?"
	err := db.QueryRow(query, geonameid).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no record exists, return 0 (or 1 for default popularity)
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
